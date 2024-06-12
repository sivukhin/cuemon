package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	cueJson "cuelang.org/go/pkg/encoding/json"

	"github.com/sivukhin/cuemon/lib/auth"
)

func ReadJson[T any](input string) (data T, source string, err error) {
	var content []byte
	if input == "stdin" {
		content, err = io.ReadAll(os.Stdin)
		source = "stdin"
	} else {
		content, err = os.ReadFile(input)
		source = fmt.Sprintf("file '%v'", input)
	}
	if err != nil {
		err = fmt.Errorf("unable to read input file '%v': %w", input, err)
		return
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal input file '%v': %w", input, err)
		return
	}
	return
}

func Export(inputs []string) error {
	buildInstances := load.Instances(inputs, &load.Config{})
	cueContext := cuecontext.New()
	values, err := cueContext.BuildInstances(buildInstances)
	if err != nil {
		return fmt.Errorf("unable to build context from files: %+v", inputs)
	}
	if len(values) != 1 {
		return fmt.Errorf("expected single value during CUE evaluation, got %v", len(values))
	}
	value := values[0]

	mainValue := value.Eval()
	rowsValue := value.LookupPath(cue.ParsePath("#rows")).Eval()
	mainJson, err := cueJson.Marshal(mainValue)
	if err != nil {
		return fmt.Errorf("unable to format main json: %v", err)
	}
	rowsJson, err := cueJson.Marshal(rowsValue)
	if err != nil {
		return fmt.Errorf("unable to format rows json: %v", err)
	}

	var mainGrafana map[string]any
	if err = json.Unmarshal([]byte(mainJson), &mainGrafana); err != nil {
		return fmt.Errorf("unable to parse main json back: %v", err)
	}

	var rowsCue []CueRow
	if err = json.Unmarshal([]byte(rowsJson), &rowsCue); err != nil {
		return fmt.Errorf("unable to parse rows json back: %v", err)
	}

	err = fillGridPositions(rowsCue)
	if err != nil {
		return fmt.Errorf("failed to fill grid positions: %v", err)
	}

	panels, err := createPanels(rowsCue)
	if err != nil {
		return fmt.Errorf("failed to create panels for final export: %w", err)
	}

	mainGrafana["panels"] = panels
	mainBytes, err := json.MarshalIndent(mainGrafana, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize final struct: %w", err)
	}

	fmt.Printf("%v", string(mainBytes))
	return nil
}

func createPanel(panel JsonRaw[GrafanaPanel]) (map[string]any, error) {
	var data map[string]any
	err := json.Unmarshal(panel.Raw, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize panel from json back: %w", err)
	}
	data["id"] = panel.Value.Id
	data["gridPos"] = map[string]int{
		"x": panel.Value.GridPos.X,
		"y": panel.Value.GridPos.Y,
		"w": panel.Value.GridPos.W,
		"h": panel.Value.GridPos.H,
	}
	return data, nil
}

func createPanels(rows []CueRow) ([]map[string]any, error) {
	panels := make([]map[string]any, 0)
	for _, row := range rows {
		rowPanels := make([]map[string]any, 0)
		for _, group := range row.Groups {
			for _, panel := range group.Panels {
				panelMap, err := createPanel(panel)
				if err != nil {
					return nil, fmt.Errorf("failed to create panel: %w", err)
				}
				rowPanels = append(rowPanels, panelMap)
			}
		}
		if row.Collapsed {
			panels = append(panels, map[string]any{
				"type":      "row",
				"id":        row.Id,
				"title":     row.Title,
				"collapsed": row.Collapsed,
				"panels":    rowPanels,
				"gridPos": map[string]int{
					"h": 1,
					"w": 24,
					"x": 0,
					"y": row.Y,
				},
			})
		} else {
			panels = append(panels, map[string]any{
				"type":      "row",
				"id":        row.Id,
				"title":     row.Title,
				"collapsed": row.Collapsed,
				"gridPos": map[string]int{
					"h": 1,
					"w": 24,
					"x": 0,
					"y": row.Y,
				},
			})
			panels = append(panels, rowPanels...)
		}
	}
	return panels, nil
}

func choose[T any](values ...*T) (T, error) {
	for i := len(values) - 1; i >= 0; i-- {
		if values[i] != nil {
			return *values[i], nil
		}
	}
	return *new(T), fmt.Errorf("all values are null")
}

func fillGridPositions(rows []CueRow) error {
	id, globalOffsetY := 0, 0
	for rowI := range rows {
		row := &rows[rowI]
		row.Y = globalOffsetY
		row.Id = id
		id++
		globalOffsetY++

		rowOffsetY := 0
		if !row.Collapsed {
			rowOffsetY = globalOffsetY
		}
		for groupI := range row.Groups {
			group := &row.Groups[groupI]
			height, err := choose(row.Height, group.Height)
			if err != nil {
				return fmt.Errorf("no height configured for group: %w", err)
			}
			widths, err := choose(row.Widths, group.Widths)
			if err != nil {
				return fmt.Errorf("no widths configured for group: %w", err)
			}

			column := 0
			groupOffsetY := rowOffsetY
			groupOffsetX := 0
			for panelI := range group.Panels {
				panel := &group.Panels[panelI]

				panel.Value.Id = id
				id++

				if groupOffsetY == rowOffsetY+height {
					groupOffsetY = rowOffsetY
					groupOffsetX += widths[column]
					column++
				}

				if column >= len(widths) {
					return fmt.Errorf("invalid panel tiling: too many rows in a group")
				}

				if panel.Value.GridPos.H == 0 {
					panel.Value.GridPos.H = (rowOffsetY + height) - groupOffsetY
				}
				panel.Value.GridPos.X = groupOffsetX
				panel.Value.GridPos.Y = groupOffsetY
				panel.Value.GridPos.W = widths[column]
				groupOffsetY += panel.Value.GridPos.H
			}

			rowOffsetY = groupOffsetY
		}

		if !row.Collapsed {
			globalOffsetY = rowOffsetY
		}
	}
	return nil
}

func Bootstrap(input, module, dir string, overwrite bool) error {
	if module == "" {
		return fmt.Errorf("module should be provided")
	}
	if dir == "" {
		return fmt.Errorf("dir should be provided")
	}
	grafana, source, err := ReadJson[Grafana](input)
	if err != nil {
		return fmt.Errorf("unable to get grafana dashboard json from %v: %w", source, err)
	}

	monitoring, err := MonitoringFiles(module, dir, grafana)
	if err != nil {
		return fmt.Errorf("unable to create monitoring files: %w", err)
	}

	log.Printf("ready to write %v files for monitoring bootstrap", len(monitoring))
	for _, file := range monitoring {
		log.Printf("  %v (%v bytes)", file.Path, len(file.Content))
	}

	errs := make([]error, 0)
	for _, file := range monitoring {
		err = os.MkdirAll(path.Dir(file.Path), os.ModePerm)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		mode := os.O_CREATE | os.O_RDWR | os.O_TRUNC
		if !overwrite {
			mode |= os.O_EXCL
		}
		f, err := os.OpenFile(file.Path, mode, os.ModePerm)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		_, err = f.Write([]byte(file.Content))
		if err != nil {
			errs = append(errs, err)
		}
		_ = f.Close()
	}
	return errors.Join(errs...)
}

func CutNode(content []byte, oldNode ast.Node, newNode ast.Node) ([]byte, []byte, []byte, error) {
	newPart, err := format.Node(newNode, format.Simplify())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to format new node: %v", err)
	}
	oldPart, err := format.Node(oldNode, format.Simplify())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to format old node: %v", err)
	}
	start, end := oldNode.Pos().Line()-1, oldNode.End().Line()
	lines := bytes.Split(content, []byte("\n"))
	var result []byte
	result = append(result, bytes.Join(lines[0:start], []byte("\n"))...)
	result = append(result, newPart...)
	result = append(result, bytes.Join(lines[end:], []byte("\n"))...)
	return result, oldPart, newPart, nil
}

func Update(input, dir string, overwrite bool) error {
	if dir == "" {
		return fmt.Errorf("update dir should be provided")
	}
	panel, source, err := ReadJson[JsonRaw[GrafanaPanel]](input)
	if err != nil {
		return fmt.Errorf("unable to get grafana panel json from %v: %w", source, err)
	}
	updates, err := UpdateFiles(dir, panel)
	if err != nil {
		return fmt.Errorf("unable to prepare updates: %w", err)
	}
	if len(updates) == 0 {
		return fmt.Errorf("unable to find configuration for panel with name '%v'", panel.Value.Title)
	}
	if len(updates) > 1 {
		var filenames []string
		for _, update := range updates {
			filenames = append(filenames, update.Path)
		}
		return fmt.Errorf("found more than 1 configuration for panel with name '%v': files=%v", panel.Value.Title, filenames)
	}
	log.Printf("--- OLD (%v) ---\n", updates[0].Path)
	log.Printf("%v\n", string(updates[0].BeforePart))
	log.Printf("--- NEW (%v) ---\n", updates[0].Path)
	log.Printf("%v\n", string(updates[0].AfterPart))
	if !overwrite {
		log.Printf("skipped update of file %v - pass overwrite flag to force this action", updates[0].Path)
		return nil
	}
	err = os.WriteFile(updates[0].Path, updates[0].Content, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("updated file %v", updates[0].Path)
	return nil
}

type DashboardPayload struct {
	Dashboard json.RawMessage `json:"dashboard"`
	Message   string          `json:"message"`
	Overwrite bool            `json:"overwrite"`
}

func searchDashboards(grafanaUrl string, authorization auth.AuthorizationMethods, name string) ([]Grafana, error) {
	request, err := http.NewRequest("GET", grafanaUrl+"/api/search", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request for url %v: %w", grafanaUrl, err)
	}
	q := request.URL.Query()
	q.Add("query", name)
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Content-Type", "application/json")
	auth.AddAuthorization(request, authorization)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send request for url %v: %w", request.URL.String(), err)
	}
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to get response for url %v: %w", request.URL.String(), err)
	}
	var dashboards []Grafana
	err = json.Unmarshal(payload, &dashboards)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response for url %v: %w", request.URL.String(), err)
	}
	return dashboards, nil
}

func updateDashboard(grafanaUrl string, authorization auth.AuthorizationMethods, payload DashboardPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to serialize request body for url %v: %w", grafanaUrl, err)
	}
	request, err := http.NewRequest("POST", grafanaUrl+"/api/dashboards/db", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("unable to create request for url %v: %w", grafanaUrl, err)
	}
	request.Header.Set("Content-Type", "application/json")
	auth.AddAuthorization(request, authorization)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("unable to send request for url %v: %w", request.URL.String(), err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read resopnse body for url %v: %w", request.URL.String(), err)
	}
	if response.StatusCode/100 != 2 {
		return fmt.Errorf("non successful HTTP status %v (body: %v)", response.StatusCode, string(body))
	}
	return nil
}

const (
	IdTag    = "Id"
	UidTag   = "Uid"
	TestTag  = "Test"
	TitleTag = "Title"
)

func Push(grafanaUrl string, authorization auth.AuthorizationMethods, dashboard string, message string, dashboardPlayground string) error {
	if grafanaUrl == "" {
		return fmt.Errorf("grafana URL should be provided")
	}
	if dashboard == "" {
		return fmt.Errorf("dashboard should be provided")
	}
	if message == "" {
		return fmt.Errorf("message should be non empty")
	}
	tags := make([]string, 0)
	if dashboardPlayground != "" {
		log.Printf("search for playground dashboard")
		dashboards, err := searchDashboards(grafanaUrl, authorization, dashboardPlayground)
		if err != nil {
			return err
		}
		if len(dashboards) == 0 {
			return fmt.Errorf("not found any dashboard matching query '%v'", dashboardPlayground)
		}
		if len(dashboards) > 1 {
			names := make([]string, 0)
			for _, d := range dashboards {
				names = append(names, d.Value.Title)
			}
			return fmt.Errorf("found many dashboards matching query '%v': %v", dashboardPlayground, strings.Join(names, ", "))
		}
		playgroundDashboard := dashboards[0]
		tags = append(tags, fmt.Sprintf("%v=%v", IdTag, playgroundDashboard.Value.Id))
		tags = append(tags, fmt.Sprintf("%v=%v", UidTag, playgroundDashboard.Value.Uid))
		tags = append(tags, fmt.Sprintf("%v=%v", TitleTag, playgroundDashboard.Value.Title))
		tags = append(tags, fmt.Sprintf("%v=%v", TestTag, true))
		log.Printf("found playground dashboard with id: %v, uid: %v, title: %v", playgroundDashboard.Value.Id, playgroundDashboard.Value.Uid, playgroundDashboard.Value.Title)
	}
	cueContext := cuecontext.New()
	dir, file := path.Split(dashboard)
	originalBuildInstances := load.Instances([]string{file}, &load.Config{
		Dir:  dir,
		Tags: tags,
	})
	original, err := cueContext.BuildInstances(originalBuildInstances)
	if err != nil {
		return fmt.Errorf("unable to load dashboard: %w", err)
	}
	grafana := original[0].LookupPath(cue.MakePath(cue.Str(GrafanaField)))
	serialized, err := cueJson.Marshal(grafana)
	if err != nil {
		return fmt.Errorf("unable to export CUE value to json: %w", err)
	}
	payload := DashboardPayload{
		Dashboard: []byte(serialized),
		Message:   message,
		Overwrite: true,
	}
	_ = payload
	log.Printf("prepared payload for dashboard update")
	err = updateDashboard(grafanaUrl, authorization, payload)
	if err != nil {
		return fmt.Errorf("unable to update grafana dashboard: %w", err)
	}
	log.Printf("successfully updated grafana dashboard")
	return nil
}
