package lib

import (
	"bytes"
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	cueJson "cuelang.org/go/pkg/encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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

func searchDashboards(grafanaUrl string, cookie string, name string) ([]Grafana, error) {
	request, err := http.NewRequest("GET", grafanaUrl+"/api/search", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request for url %v: %w", grafanaUrl, err)
	}
	q := request.URL.Query()
	q.Add("query", name)
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", cookie)
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

func updateDashboard(grafanaUrl string, cookie string, payload DashboardPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to serialize request body for url %v: %w", grafanaUrl, err)
	}
	request, err := http.NewRequest("POST", grafanaUrl+"/api/dashboards/db", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("unable to create request for url %v: %w", grafanaUrl, err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", cookie)
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

func Push(grafanaUrl string, cookie string, dashboard string, message string, dashboardTmp string) error {
	if grafanaUrl == "" {
		return fmt.Errorf("grafana URL should be provided")
	}
	if dashboard == "" {
		return fmt.Errorf("dashboard should be provided")
	}
	if message == "" {
		return fmt.Errorf("message should be non empty")
	}
	if cookie == "" {
		return fmt.Errorf("cookie should be provided implicitly through RunContext in interactive mode or explicitly through GRAFANA_COOKIE env variable")
	}
	tags := make([]string, 0)
	if dashboardTmp != "" {
		log.Printf("search for temp dashboard")
		dashboards, err := searchDashboards(grafanaUrl, cookie, dashboardTmp)
		if err != nil {
			return err
		}
		if len(dashboards) == 0 {
			return fmt.Errorf("not found any dashboard matching query '%v'", dashboardTmp)
		}
		if len(dashboards) > 1 {
			names := make([]string, 0)
			for _, d := range dashboards {
				names = append(names, d.Value.Title)
			}
			return fmt.Errorf("found many dashboards matching query '%v': %v", dashboardTmp, strings.Join(names, ", "))
		}
		tempDashboard := dashboards[0]
		tags = append(tags, fmt.Sprintf("%v=%v", IdTag, tempDashboard.Value.Id))
		tags = append(tags, fmt.Sprintf("%v=%v", UidTag, tempDashboard.Value.Uid))
		tags = append(tags, fmt.Sprintf("%v=%v", TitleTag, tempDashboard.Value.Title))
		tags = append(tags, fmt.Sprintf("%v=%v", TestTag, true))
		log.Printf("found temp dashboard with id: %v, uid: %v, title: %v", tempDashboard.Value.Id, tempDashboard.Value.Uid, tempDashboard.Value.Title)
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
	log.Printf("prepared payload for dashboard update")
	err = updateDashboard(grafanaUrl, cookie, payload)
	if err != nil {
		return fmt.Errorf("unable to update grafana dashboard: %w", err)
	}
	log.Printf("successfully updated grafana dashboard")
	return nil
}
