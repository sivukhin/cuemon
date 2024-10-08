package lib

import (
	_ "embed"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/token"
)

var (
	//go:embed cue/grafanaV10.cue
	GrafanaV10Cue string

	//go:embed cue/monMeta.cue
	MonMetaCue string
	//go:embed cuebootstrap/monMetaBootstrap.cue
	MonMetaBootstrapCue string

	//go:embed cue/monLink.cue
	MonLinkCue string
	//go:embed cuebootstrap/monLinkBootstrap.cue
	MonLinkBootstrapCue string

	//go:embed cue/monVariable.cue
	MonVariableCue string
	//go:embed cuebootstrap/monVariableBootstrap.cue
	MonVariableBootstrapCue string

	//go:embed cue/monPanel.cue
	MonPanelCue string
	//go:embed cuebootstrap/monPanelBootstrap.cue
	MonPanelBootstrapCue string
)

var matchIgnoredSymbols = regexp.MustCompile("[^a-zA-Z0-9]")
var matchUnderscores = regexp.MustCompile("_+")

func toFilename(s string) string {
	s = matchIgnoredSymbols.ReplaceAllString(s, "_")
	s = matchUnderscores.ReplaceAllString(s, "_")
	return strings.ToLower(s)
}

const (
	TitleField     = "Title"
	CollapsedField = "Collapsed"
	ColumnsField   = "Widths"
	HeightsField   = "Heights"
	PanelField     = "Panel"
	PanelGridField = "PanelGrid"
	WidthField     = "Width"
	HeightField    = "Height"
)

type MonitoringContext struct {
	SchemaFiles    []string
	ReductionRules ReductionRules
}

type CuemonPanel struct {
	GridPos Box
	Src     []ast.Decl
}

type CuemonRow struct {
	Title     string
	Collapsed bool
	Panels    []CuemonPanel
}

func (m MonitoringContext) cueConvert(
	conversions []string,
	context map[string]string,
	trimAgainstVar string,
	mappingName string,
	fieldName string,
	flatten bool,
) ([]ast.Decl, error) {
	decls, err := CueConvert(conversions, context, trimAgainstVar, mappingName, fieldName, flatten)
	if err != nil {
		return nil, err
	}
	return m.ReductionRules.ReduceAst(decls), nil
}

func (m MonitoringContext) CreateMetaSrc(grafana Grafana) ([]ast.Decl, error) {
	schemas := append(m.SchemaFiles, MonMetaCue, MonMetaBootstrapCue)
	context := map[string]string{"input": string(grafana.Raw)}
	return m.cueConvert(schemas, context, "#monMeta", "#conversion", "output", false)
}

func (m MonitoringContext) CreateLinkSrc(link Link) ([]ast.Decl, error) {
	schemas := append(m.SchemaFiles, MonLinkCue, MonLinkBootstrapCue)
	context := map[string]string{"input": string(link.Raw)}
	return m.cueConvert(schemas, context, "#monLink", "#conversion", "output", false)
}

func (m MonitoringContext) CreateVariableSrc(variable Templating) ([]ast.Decl, error) {
	schemas := append(m.SchemaFiles, MonVariableCue, MonVariableBootstrapCue)
	context := map[string]string{"input": string(variable.Raw)}
	return m.cueConvert(schemas, context, "#monVariable", "#conversion", "output", false)
}

func (m MonitoringContext) CreatePanelSrc(panel JsonRaw[GrafanaPanel]) ([]ast.Decl, error) {
	schemas := append(m.SchemaFiles, MonPanelCue, MonPanelBootstrapCue)
	context := map[string]string{"input": string(panel.Raw)}
	return m.cueConvert(schemas, context, "#monPanel", "#conversion", "output", false)
}

func (m MonitoringContext) CreateRowSrc(row CuemonRow) (string, []ast.Decl, error) {
	grid := make([]Box, 0)
	for i, panel := range row.Panels {
		box := panel.GridPos
		box.Id = i
		grid = append(grid, box)
	}
	gridLayout := AlignGrid(grid)

	var rowDecls []ast.Decl
	rowDecls = append(rowDecls, FieldIdent("title", ast.NewString(row.Title)))
	rowDecls = append(rowDecls, FieldIdent("collapsed", ast.NewBool(row.Collapsed)))
	rowDecls = append(rowDecls, FieldIdent("h", Int(gridLayout.Height)))
	rowDecls = append(rowDecls, FieldIdent("w", IntList(gridLayout.Widths)))
	rowDecls = append(rowDecls, LineBreak())

	var groupsDecls []ast.Expr
	for _, group := range gridLayout.Groups {
		var groupDecls []ast.Decl
		groupDecls = append(groupDecls, FieldIdent("h", Int(group.Height)))
		groupDecls = append(groupDecls, FieldIdent("w", IntList(group.Widths)))

		var panelsDecls []ast.Expr
		for _, panel := range group.Panels {
			var panelDecls []ast.Decl

			if override, ok := group.Overrides[panel.Id]; ok {
				if override.Height > 0 {
					panelDecls = append(panelDecls, FieldIdent("h", Int(override.Height)))
				}
			}
			panelDecls = append(panelDecls, row.Panels[panel.Id].Src...)
			panelsDecls = append(panelsDecls, ast.NewStruct(AsAny(panelDecls)...))
		}

		groupDecls = append(groupDecls, FieldIdent("panels", ast.NewList(panelsDecls...)))
		groupsDecls = append(groupsDecls, ast.NewStruct(AsAny(groupDecls)...))
	}
	rowDecls = append(rowDecls, FieldIdent("groups", ast.NewList(groupsDecls...)))
	return row.Title, rowDecls, nil
}

func (m MonitoringContext) createRow(row GrafanaPanel) (*Row, error) {
	decls := []ast.Decl{FieldIdent(TitleField, ast.NewString(row.Title))}
	if row.Collapsed {
		decls = append(decls, FieldIdent(CollapsedField, ast.NewBool(true)))
	}

	grid := make([]Box, 0)
	for i, panel := range row.Panels {
		box := panel.Value.GridPos
		box.Id = i
		grid = append(grid, box)
	}
	layout := AnalyzeGrid(grid)

	decls = append(decls, FieldIdent(ColumnsField, IntList(layout.Columns)))
	if len(layout.Heights) == 1 {
		decls = append(decls, FieldIdent(HeightsField, Int(layout.Heights[0])))
	} else {
		decls = append(decls, FieldIdent(HeightsField, IntList(layout.Heights)))
	}
	decls = append(decls, LineBreak())

	for _, id := range layout.Order {
		panel := row.Panels[id]
		panelDecls, err := m.CreatePanelSrc(panel)
		if err != nil {
			return nil, fmt.Errorf("unable to convert panel %v: %w", panel.Value.Title, err)
		}
		if override, ok := layout.Overrides[id]; ok {
			fields := []any{FieldIdent(WidthField, Int(override.Width))}
			if override.Height > 0 {
				fields = append(fields, FieldIdent(HeightField, Int(override.Height)))
			}
			decls = append(decls, FieldIdent(PanelGridField, ast.NewStruct(ast.NewIdent(panel.Value.Title), ast.NewStruct(fields...))))
		}
		decls = append(decls, FieldIdent(PanelField, ast.NewStruct(
			ast.NewIdent(panel.Value.Title), ast.NewStruct(AsAny(panelDecls)...),
		)))
		decls = append(decls, LineBreak())
	}
	return &Row{Name: toFilename(row.Title), Src: decls}, nil
}

type Row struct {
	Name string
	Src  []ast.Decl
}

type Monitoring struct {
	Meta      []ast.Decl
	Variables []ast.Expr
	Links     []ast.Expr
	RowNames  []string
	Rows      []ast.Expr
}

func (m MonitoringContext) ExtractRows(grafana Grafana) []CuemonRow {
	title, collapsed := "top_unnamed", false
	i := 0
	rows := make([]CuemonRow, 0)
	for i < len(grafana.Value.Panels) {
		current := grafana.Value.Panels[i]
		if current.Value.Type == "row" {
			title = current.Value.Title
			collapsed = current.Value.Collapsed
			i += 1
		}
		var panels []JsonRaw[GrafanaPanel]
		if collapsed {
			panels = current.Value.Panels
		} else {
			startI := i
			for i < len(grafana.Value.Panels) && grafana.Value.Panels[i].Value.Type != "row" {
				i += 1
			}
			panels = grafana.Value.Panels[startI:i]
		}
		cuemonPanels := make([]CuemonPanel, 0, len(panels))
		for panelI, panel := range panels {
			src, err := m.CreatePanelSrc(panel)
			if err != nil {
				Logger.Fatalf("failed to create panel source (title: %v): %v", panel.Value.Title, err)
			}
			box := panel.Value.GridPos
			box.Id = panelI
			cuemonPanels = append(cuemonPanels, CuemonPanel{GridPos: box, Src: src})
		}
		rows = append(rows, CuemonRow{Title: title, Collapsed: collapsed, Panels: cuemonPanels})
	}
	return rows
}

func (m MonitoringContext) CreateDashboardSrc(grafana Grafana) (*Monitoring, error) {
	meta, err := m.CreateMetaSrc(grafana)

	if err != nil {
		return nil, err
	}
	links := make([]ast.Expr, 0)
	for _, linkGrafana := range grafana.Value.Links {
		log.Printf("bootstraping CUE node for link...")
		link, err := m.CreateLinkSrc(linkGrafana)

		if err != nil {
			return nil, err
		}
		links = append(links, ast.NewStruct(AsAny(link)...))
	}
	variables := make([]ast.Expr, 0)
	for _, template := range grafana.Value.Templating.List {
		log.Printf("bootstraping CUE node for '%v' variable...", template.Value.Name)
		variable, err := m.CreateVariableSrc(template)
		if err != nil {
			return nil, err
		}
		variables = append(variables, ast.NewStruct(AsAny(variable)...))
	}
	rows := make([]ast.Expr, 0)
	rowNames := make([]string, 0)
	for _, grafanaRow := range m.ExtractRows(grafana) {
		log.Printf("bootstraping CUE file for '%v' row...", grafanaRow.Title)
		name, row, err := m.CreateRowSrc(grafanaRow)
		if err != nil {
			return nil, err
		}
		rows = append(rows, ast.NewStruct(AsAny(row)...))
		rowNames = append(rowNames, name)
	}
	monitoring := &Monitoring{
		Meta:      meta,
		Links:     links,
		Variables: variables,
		Rows:      rows,
		RowNames:  rowNames,
	}
	return monitoring, nil
}

type FileEntry struct {
	Path    string
	Content string
}

func MonitoringFiles(module, output string, grafana Grafana) ([]FileEntry, error) {
	dataSources, dataSourcesId := make(map[string]ast.Expr), make(map[string]string)
	projectNames, projectNamesId := make(map[string]ast.Expr), make(map[string]string)
	var grafanaSchema string
	if grafana.Value.SchemaVersion >= 39 {
		grafanaSchema = GrafanaV10Cue
	} else {
		return nil, fmt.Errorf("unsupported Grafana schema version: %v", grafana.Value.SchemaVersion)
	}
	m := MonitoringContext{
		SchemaFiles: []string{grafanaSchema},
		ReductionRules: []ReductionRule{
			{
				Reduction: func(path []string, value ast.Expr) (string, ast.Expr, bool) {
					if strings.HasPrefix(path[len(path)-1], "#") {
						return "", nil, false
					}
					if structLit, ok := value.(*ast.StructLit); ok && len(structLit.Elts) == 0 {
						return "", nil, true
					}
					if arrayLit, ok := value.(*ast.ListLit); ok {
						shouldDelete := true
						for _, decl := range arrayLit.Elts {
							if structLit, ok := decl.(*ast.StructLit); !ok || len(structLit.Elts) > 0 {
								shouldDelete = false
							}
						}
						if shouldDelete {
							return "", nil, true
						}
					}
					if basicLit, ok := value.(*ast.BasicLit); ok && basicLit.Value == "null" {
						return "", nil, true
					}
					return path[len(path)-1], value, false
				},
			},
			{
				Reduction: func(path []string, expr ast.Expr) (string, ast.Expr, bool) {
					if path[len(path)-1] != "#datasrc" {
						return "", nil, false
					}
					dataSourceSrc, err := FormatNode(expr)
					if err != nil {
						Logger.Errorf("failed to format data source: %v", err)
					}
					if _, ok := dataSourcesId[dataSourceSrc]; !ok {
						dataSourceId := fmt.Sprintf("#dataSource%v", len(dataSources)+1)
						dataSourcesId[dataSourceSrc] = dataSourceId
						dataSources[dataSourceId] = expr
					}
					return path[len(path)-1], ast.NewIdent(dataSourcesId[dataSourceSrc]), true
				},
			},
			{
				Reduction: func(path []string, expr ast.Expr) (string, ast.Expr, bool) {
					if path[len(path)-1] != "projectName" {
						return "", nil, false
					}
					projectNameSrc, err := FormatNode(expr)
					if err != nil {
						Logger.Errorf("failed to format project name: %v", err)
					}
					if _, ok := projectNamesId[projectNameSrc]; !ok {
						projectNameId := fmt.Sprintf("#projectName%v", len(projectNames)+1)
						projectNamesId[projectNameSrc] = projectNameId
						projectNames[projectNameId] = expr
					}
					return path[len(path)-1], ast.NewIdent(projectNamesId[projectNameSrc]), true
				},
			},
			{
				Reduction: func(path []string, expr ast.Expr) (string, ast.Expr, bool) {
					if path[0] == "version" || path[0] == "iteration" || path[0] == "id" || path[0] == "gridPos" {
						return "", nil, true
					}
					if path[len(path)-1] == "hashKey" {
						return "", nil, true
					}
					return "", nil, false
				},
			},
		},
	}
	monitoring, err := m.CreateDashboardSrc(grafana)
	if err != nil {
		return nil, fmt.Errorf("unable to create monitoring configuration: %w", err)
	}
	files := make([]FileEntry, 0)

	variablesSrc, err := FormatDecls([]ast.Decl{
		Package(module),
		FieldIdent("#variables", ast.NewList(monitoring.Variables...)),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to format variable source: %w", err)
	}
	rowsRefs := make([]ast.Expr, 0)
	files = append(files, FileEntry{path.Join(output, "variables.cue"), variablesSrc})
	for rowI, row := range monitoring.Rows {
		rowName := CapitalizeName("row_" + monitoring.RowNames[rowI])
		rowSrc, err := FormatDecls([]ast.Decl{
			Package(module),
			FieldIdent("#"+rowName, ast.NewBinExpr(token.AND, ast.NewIdent("#monRow"), row)),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to format row '%v' source: %w", rowName, err)
		}
		rowsRefs = append(rowsRefs, ast.NewIdent("#"+rowName))
		rowFileName := fmt.Sprintf("%v.cue", rowName)
		files = append(files, FileEntry{path.Join(output, rowFileName), rowSrc})
	}

	dashboardDecls := []ast.Decl{
		Package(module),
		ast.NewIdent("#mon"),
	}
	dashboardDecls = append(dashboardDecls, monitoring.Meta...)

	for dataSourceName, dataSourceExpr := range dataSources {
		dashboardDecls = append(dashboardDecls, FieldIdent(dataSourceName, dataSourceExpr))
	}

	for projectName, projectNameExpr := range projectNames {
		dashboardDecls = append(dashboardDecls, FieldIdent(projectName, projectNameExpr))
	}

	dashboardDecls = append(dashboardDecls, FieldIdent("#links", ast.NewList(monitoring.Links...)))
	dashboardDecls = append(dashboardDecls, FieldIdent("#rows", ast.NewList(rowsRefs...)))
	dashboardSrc, err := FormatDecls(dashboardDecls)
	if err != nil {
		return nil, fmt.Errorf("unable to format dashboard source: %w", err)
	}
	files = append(files, FileEntry{path.Join(output, "dashboard.cue"), dashboardSrc})
	return files, nil
}
