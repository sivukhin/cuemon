package lib

import (
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/token"
)

var (
	//go:embed cue/grafana_v7.cue
	GrafanaV7Cue string
	//go:embed cue/grafana_v10.cue
	GrafanaV10Cue string

	//go:embed cue/mon_meta.cue
	MonMetaCue string
	//go:embed cue/mon_meta_bootstrap.cue
	MonMetaBootstrapCue string

	//go:embed cue/mon_link.cue
	MonLinkCue string
	//go:embed cue/mon_link_bootstrap.cue
	MonLinkBootstrapCue string

	//go:embed cue/mon_variable.cue
	MonVariableCue string
	//go:embed cue/mon_variable_bootstrap.cue
	MonVariableBootstrapCue string

	//go:embed cue/mon_panel.cue
	MonPanelCue string
	//go:embed cue/mon_panel_bootstrap.cue
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
	GrafanaField       = "Grafana"
	SchemaVersionField = "schemaVersion"
	TitleField         = "Title"
	CollapsedField     = "Collapsed"
	ColumnsField       = "Widths"
	HeightsField       = "Heights"
	PanelField         = "Panel"
	PanelGridField     = "PanelGrid"
	WidthField         = "Width"
	HeightField        = "Height"
)

type MonitoringContext struct {
	SchemaFiles    []string
	ReductionRules ReductionRules
}

//func (m MonitoringContext) matchOverrides(legend string, overrides []GrafanaSeriesOverrides) ([]ast.Decl, bool) {
//	for _, override := range overrides {
//		if strings.HasPrefix(override.Alias, "/") && strings.HasSuffix(override.Alias, "/") {
//			matched, err := regexp.MatchString(strings.Trim(override.Alias, "/"), legend)
//			if err != nil {
//				log.Printf("unable to match legend %v against alias %v: %v", legend, override.Alias, err)
//				continue
//			}
//			if !matched {
//				continue
//			}
//		} else {
//			if override.Alias != legend {
//				continue
//			}
//		}
//
//		fields := make([]any, 0)
//		if override.Dashes != nil {
//			fields = append(fields, ast.NewIdent(DashesField), ast.NewBool(*override.Dashes))
//		}
//		if override.Hidden != nil {
//			fields = append(fields, ast.NewIdent(HiddenField), ast.NewBool(*override.Hidden))
//		}
//		if override.Color != nil {
//			fields = append(fields, ast.NewIdent(ColorField), ast.NewString(*override.Color))
//		}
//		if override.Fill != nil {
//			fields = append(fields, ast.NewIdent(FillField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.Fill)})
//		}
//		if override.YAxis != nil {
//			fields = append(fields, ast.NewIdent(YAxisField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.YAxis)})
//		}
//		if override.ZIndex != nil {
//			fields = append(fields, ast.NewIdent(ZIndexField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.ZIndex)})
//		}
//		if override.LineWidth != nil {
//			fields = append(fields, ast.NewIdent(LineWidthField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.LineWidth)})
//		}
//		return []ast.Decl{FieldIdent(OverridesField, ast.NewStruct(fields...))}, true
//	}
//	return nil, false
//}

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
	dataSources := make(map[string]ast.Expr)
	dataSourcesId := make(map[string]string)
	m := MonitoringContext{
		SchemaFiles: []string{GrafanaV7Cue},
		ReductionRules: []ReductionRule{
			{
				Reduction: func(path []string, value ast.Expr) (string, ast.Expr, bool) {
					if structLit, ok := value.(*ast.StructLit); ok && len(structLit.Elts) == 0 {
						return "", nil, true
					}
					if basicLit, ok := value.(*ast.BasicLit); ok && basicLit.Value == "null" {
						return "", nil, true
					}
					return path[len(path)-1], value, false
				},
			},
			{
				Reduction: func(path []string, expr ast.Expr) (string, ast.Expr, bool) {
					if path[len(path)-1] != "datasource" {
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
					if path[0] == "version" || path[0] == "iteration" || path[0] == "id" || path[0] == "gridPos" {
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

	dashboardDecls = append(dashboardDecls, FieldIdent("#links", ast.NewList(monitoring.Links...)))
	dashboardDecls = append(dashboardDecls, FieldIdent("#rows", ast.NewList(rowsRefs...)))
	dashboardSrc, err := FormatDecls(dashboardDecls)
	if err != nil {
		return nil, fmt.Errorf("unable to format dashboard source: %w", err)
	}
	files = append(files, FileEntry{path.Join(output, "dashboard.cue"), dashboardSrc})
	return files, nil
}

type FileEntryUpdated struct {
	Path       string
	Content    []byte
	BeforePart []byte
	AfterPart  []byte
}

func UpdateFiles(output string, panel JsonRaw[GrafanaPanel]) ([]FileEntryUpdated, error) {
	targets := make(map[string]ast.Node, 0)
	versions := make([]int, 0)
	definitions := make(map[string]*ast.BasicLit)
	_ = filepath.WalkDir(output, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return fmt.Errorf("directory %v not found", output)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".cue") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		cueAst, err := CueAst(string(content))
		if err != nil {
			return err
		}
		for _, decl := range cueAst.Decls {
			field, ok := decl.(*ast.Field)
			if !ok {
				continue
			}
			ident, ok := field.Label.(*ast.Ident)
			if !ok {
				continue
			}
			if !strings.HasPrefix(ident.Name, "#") {
				continue
			}
			value, ok := field.Value.(*ast.BasicLit)
			if !ok {
				continue
			}
			definitions[ident.Name] = value
		}

		_, versionNode := findField(cueAst.Decls, GrafanaField, SchemaVersionField)
		if versionNode != nil {
			version, err := strconv.Atoi(versionNode.(*ast.BasicLit).Value)
			if err != nil {
				return fmt.Errorf("unable to parse Grafana schema version: %w", err)
			}
			versions = append(versions, version)
		}
		target, _ := findField(cueAst.Decls, PanelField, panel.Value.Title)
		if target != nil {
			targets[path] = target
		}
		return nil
	})
	if len(Unique(versions)) != 1 {
		return nil, fmt.Errorf("unable to find single Grafana schema version: %v", versions)
	}
	//version := versions[0]
	m := MonitoringContext{}
	panelDecls, err := m.CreatePanelSrc(panel)
	if err != nil {
		return nil, fmt.Errorf("unable to convert panel %v: %w", panel.Value.Title, err)
	}
	panelStruct := FieldIdent(PanelField, ast.NewStruct(
		ast.NewIdent(panel.Value.Title), ast.NewStruct(AsAny(panelDecls)...),
	))
	astutil.Apply(panelStruct, nil, func(cursor astutil.Cursor) bool {
		lit, ok := cursor.Node().(*ast.BasicLit)
		if !ok {
			return true
		}
		for definition, replacement := range definitions {
			if replacement.Value == lit.Value && replacement.Kind == lit.Kind {
				cursor.Replace(ast.NewIdent(definition))
				return true
			}
		}
		return true
	})
	updates := make([]FileEntryUpdated, 0)
	for targetPath, target := range targets {
		content, err := os.ReadFile(targetPath)
		if err != nil {
			return nil, err
		}
		cutted, oldPart, newPart, err := CutNode(content, target, panelStruct)
		if err != nil {
			return nil, err
		}
		updates = append(updates, FileEntryUpdated{Path: targetPath, Content: cutted, BeforePart: oldPart, AfterPart: newPart})
	}
	return updates, nil
}
