package lib

import (
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/literal"
	"cuelang.org/go/cue/token"
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
)

var (
	//go:embed cue/grafana.cue
	GrafanaCue string
	//go:embed cue/mon.cue
	MonCue string
	//go:embed cue/mon2grafana.cue
	Mon2GrafanaCue string
	//go:embed grafana2mon/meta.cue
	Meta2MonitoringCue string
	//go:embed grafana2mon/template.cue
	Template2MonitoringCue string
	//go:embed grafana2mon/panel.cue
	Panel2MonitoringCue string
	//go:embed grafana2mon/target.cue
	Target2MonitoringCue string
)

var matchIgnoredSymbols = regexp.MustCompile("[^a-zA-Z0-9]")
var matchUnderscores = regexp.MustCompile("_+")

func toFilename(s string) string {
	s = matchIgnoredSymbols.ReplaceAllString(s, "_")
	s = matchUnderscores.ReplaceAllString(s, "_")
	return strings.ToLower(s)
}

const (
	GrafanaField        = "Grafana"
	SchemaVersionField  = "schemaVersion"
	TitleField          = "Title"
	CollapsedField      = "Collapsed"
	ColumnsField        = "Columns"
	HeightsField        = "Heights"
	PanelField          = "Panel"
	PanelGridField      = "PanelGrid"
	MetricsField        = "Metrics"
	WidthField          = "Width"
	HeightField         = "Height"
	VariablesField      = "Variables"
	RowsField           = "Rows"
	LegendField         = "Legend"
	OverridesField      = "Overrides"
	GrafanaIdField      = "id"
	GrafanaUidField     = "uid"
	GrafanaTitleField   = "title"
	GrafanaVersionField = "#version"
	TestField           = "Test"

	DashesField    = "Dashes"
	HiddenField    = "Hidden"
	FillField      = "Fill"
	YAxisField     = "YAxis"
	ZIndexField    = "ZIndex"
	LineWidthField = "LineWidth"
	ColorField     = "Color"
)

const (
	PanelVar     = "#Panel"
	TargetVar    = "#Target"
	RowVar       = "#Row"
	variablesVar = "variables"
	anyVar       = "_"
)

const (
	CuemonPath = "github.com/sivukhin/cuemon"
	CuemonName = "cuemon"
)

type monitoringContext struct {
	context string
}

func (m monitoringContext) matchOverrides(legend string, overrides []GrafanaSeriesOverrides) ([]ast.Decl, bool) {
	for _, override := range overrides {
		matched, err := regexp.MatchString(strings.Trim(override.Alias, "/"), legend)
		if err != nil {
			log.Printf("unable to match legend %v against alias %v: %v", legend, override.Alias, err)
			continue
		}
		if !matched {
			continue
		}
		fields := make([]any, 0)
		if override.Dashes != nil {
			fields = append(fields, ast.NewIdent(DashesField), ast.NewBool(*override.Dashes))
		}
		if override.Hidden != nil {
			fields = append(fields, ast.NewIdent(HiddenField), ast.NewBool(*override.Hidden))
		}
		if override.Color != nil {
			fields = append(fields, ast.NewIdent(ColorField), ast.NewString(*override.Color))
		}
		if override.Fill != nil {
			fields = append(fields, ast.NewIdent(FillField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.Fill)})
		}
		if override.YAxis != nil {
			fields = append(fields, ast.NewIdent(YAxisField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.YAxis)})
		}
		if override.ZIndex != nil {
			fields = append(fields, ast.NewIdent(ZIndexField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.ZIndex)})
		}
		if override.LineWidth != nil {
			fields = append(fields, ast.NewIdent(LineWidthField), &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(*override.LineWidth)})
		}
		return []ast.Decl{FieldIdent(OverridesField, ast.NewStruct(fields...))}, true
	}
	return nil, false
}

func (m monitoringContext) createPanel(panel JsonRaw[GrafanaPanel]) ([]ast.Decl, error) {
	panelDecls, err := CueConvert(PanelVar, []string{GrafanaCue, MonCue, Panel2MonitoringCue}, map[string]string{
		"Input":         string(panel.Raw),
		"SchemaVersion": m.context,
	}, false)
	if err != nil {
		return nil, fmt.Errorf("unable to convert panel '%v': %w", panel.Value.Title, err)
	}
	targetsDecls := make([]ast.Expr, 0)
	for i, target := range panel.Value.Targets {
		targetDecls, err := CueConvert(TargetVar, []string{GrafanaCue, MonCue, Target2MonitoringCue}, map[string]string{
			"Input":         string(target),
			"SchemaVersion": m.context,
		}, false)
		if err != nil {
			return nil, fmt.Errorf("unable to convert target #%v from panel '%v': %w", i, panel.Value.Title, err)
		}
		_, legendNode := findField(targetDecls, LegendField)
		if legendNode == nil {
			return nil, fmt.Errorf("legend is absent for target %v in panel %v", i, panel.Value.Title)
		}
		legend, err := literal.Unquote(legendNode.(*ast.BasicLit).Value)
		if err != nil {
			return nil, fmt.Errorf("legend for target %v in panel %v has unexpected format: %w", i, panel.Value.Title, err)
		}
		overrides, ok := m.matchOverrides(legend, panel.Value.SeriesOverrides)
		if ok {
			targetDecls = append(targetDecls, overrides...)
		}
		astutil.Apply(File(targetDecls), nil, func(cursor astutil.Cursor) bool {
			lit, ok := cursor.Node().(*ast.BasicLit)
			if !ok {
				return true
			}
			if lit.Kind != token.STRING {
				return true
			}
			value, err := literal.Unquote(lit.Value)
			if err != nil {
				log.Printf("unable to unquote string literal: %v", err)
				return true
			}
			if strings.Contains(value, "\"") && !strings.Contains(value, "#") && !strings.Contains(value, "\n") {
				cursor.Replace(&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("#\"%v\"#", value)})
			}
			return true
		})
		targetsDecls = append(targetsDecls, ast.NewStruct(AsAny(targetDecls)...))
	}
	if len(targetsDecls) > 0 {
		panelDecls = append(panelDecls, FieldIdent(MetricsField, ast.NewList(targetsDecls...)))
	}
	return panelDecls, nil
}

func (m monitoringContext) createRow(row GrafanaPanel, panels []JsonRaw[GrafanaPanel]) (*Row, error) {
	decls := []ast.Decl{FieldIdent(TitleField, ast.NewString(row.Title))}
	if row.Collapsed {
		decls = append(decls, FieldIdent(CollapsedField, ast.NewBool(true)))
	}

	grid := make([]Box, 0)
	for i, panel := range panels {
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
		panel := panels[id]
		panelDecls, err := m.createPanel(panel)
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
	Variables []ast.Decl
	Rows      []*Row
}

func (m monitoringContext) createVariable(template Templating) ([]ast.Decl, error) {
	variable, err := CueConvert(anyVar, []string{GrafanaCue, MonCue, Template2MonitoringCue}, map[string]string{
		"Input":         string(template.Raw),
		"SchemaVersion": m.context,
	}, false)
	if err != nil {
		return nil, fmt.Errorf("unable to convert variable '%v': %v", template.Value.Name, err)
	}
	return []ast.Decl{FieldIdent(template.Value.Name, ast.NewStruct(AsAny(variable)...))}, nil
}

func (m monitoringContext) creteMeta(grafana Grafana) ([]ast.Decl, error) {
	meta, err := CueConvert(anyVar, []string{GrafanaCue, Meta2MonitoringCue}, map[string]string{
		"Input":         string(grafana.Raw),
		"SchemaVersion": m.context,
	}, true)
	if err != nil {
		return nil, fmt.Errorf("unable to convert meta: %w", err)
	}
	return append([]ast.Decl{
		FieldIdent(TestField, ast.NewBinExpr(token.OR, ast.NewIdent("bool"),
			&ast.UnaryExpr{Op: token.MUL, X: ast.NewBool(false)},
		), &ast.Attribute{Text: "@tag(Test,type=bool)"}),
		FieldIdent(GrafanaField, ast.NewStruct(
			FieldIdent(GrafanaIdField, ast.NewBinExpr(token.OR, ast.NewIdent("number"),
				&ast.UnaryExpr{Op: token.MUL, X: ast.NewLit(token.INT, fmt.Sprintf("%v", grafana.Value.Id))},
			), &ast.Attribute{Text: "@tag(Id,type=number)"}),
			FieldIdent(GrafanaUidField, ast.NewBinExpr(token.OR, ast.NewIdent("string"),
				&ast.UnaryExpr{Op: token.MUL, X: ast.NewLit(token.STRING, fmt.Sprintf("\"%v\"", grafana.Value.Uid))},
			), &ast.Attribute{Text: "@tag(Uid,type=string)"}),
			FieldIdent(GrafanaTitleField, ast.NewBinExpr(token.OR, ast.NewIdent("string"),
				&ast.UnaryExpr{Op: token.MUL, X: ast.NewLit(token.STRING, fmt.Sprintf("\"%v\"", grafana.Value.Title))},
			), &ast.Attribute{Text: "@tag(Title,type=string)"}),
			FieldIdent(GrafanaVersionField, ast.NewBinExpr(token.OR, ast.NewIdent("number"),
				&ast.UnaryExpr{Op: token.MUL, X: ast.NewNull()},
			), &ast.Attribute{Text: "@tag(Version,type=number)"}),
		)),
	}, meta...), nil
}

func ExtractRows(grafana Grafana) []GrafanaPanel {
	title, collapsed := "top_unnamed", false
	i := 0
	rows := make([]GrafanaPanel, 0)
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
		rows = append(rows, GrafanaPanel{Title: title, Collapsed: collapsed, Panels: panels})
	}
	return rows
}

func (m monitoringContext) createDashboard(grafana Grafana) (*Monitoring, error) {
	meta, err := m.creteMeta(grafana)
	if err != nil {
		return nil, err
	}
	variables := make([]ast.Decl, 0)
	for _, template := range grafana.Value.Templating.List {
		log.Printf("bootstraping CUE node for '%v' variable...", template.Value.Name)
		variable, err := m.createVariable(template)
		if err != nil {
			return nil, err
		}
		variables = append(variables, variable...)
	}
	rows := make([]*Row, 0)
	for _, grafanaRow := range ExtractRows(grafana) {
		log.Printf("bootstraping CUE file for '%v' row...", grafanaRow.Title)
		row, err := m.createRow(grafanaRow, grafanaRow.Panels)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return &Monitoring{Meta: meta, Variables: variables, Rows: rows}, nil
}

type FileEntry struct {
	Path    string
	Content string
}

func MonitoringFiles(module, output string, grafana Grafana) ([]FileEntry, error) {
	m := monitoringContext{context: fmt.Sprintf(`%v`, grafana.Value.SchemaVersion)}
	monitoring, err := m.createDashboard(grafana)
	if err != nil {
		return nil, fmt.Errorf("unable to create monitoring configuration: %w", err)
	}
	packageName := PackageName(module)
	files := make([]FileEntry, 0)
	files = append(files, FileEntry{path.Join(output, "cue.mod", "module.cue"), fmt.Sprintf("module: \"%v\"\n", module)})
	files = append(files, FileEntry{path.Join(output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "grafana.cue"), GrafanaCue})
	files = append(files, FileEntry{path.Join(output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "mon.cue"), MonCue})
	files = append(files, FileEntry{path.Join(output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "mon2grafana.cue"), Mon2GrafanaCue})

	rowVars := make([]string, 0)
	dashboardImports := []string{CuemonPath, fmt.Sprintf("%v:%v", module, variablesVar)}
	for _, row := range monitoring.Rows {
		dashboardImports = append(dashboardImports, fmt.Sprintf("%v/rows:%v", module, row.Name))
		rowVars = append(rowVars, row.Name)
	}
	dashboardDecls := []ast.Decl{
		Package(packageName),
		Imports(dashboardImports...),
		ast.NewIdent(CuemonName),
		LineBreak(),
	}
	dashboardDecls = append(dashboardDecls, monitoring.Meta...)
	dashboardDecls = append(dashboardDecls, FieldIdent(VariablesField, ast.NewIdent(variablesVar)))
	dashboardDecls = append(dashboardDecls, FieldIdent(RowsField, IdentList(rowVars)))

	dashboardSrc, err := FormatDecls(dashboardDecls)
	if err != nil {
		return nil, fmt.Errorf("unable to format dashboard source: %w", err)
	}
	files = append(files, FileEntry{path.Join(output, "dashboard.cue"), dashboardSrc})

	variablesSrc, err := FormatDecls(append([]ast.Decl{Package(variablesVar)}, monitoring.Variables...))
	if err != nil {
		return nil, fmt.Errorf("unable to format variable source: %w", err)
	}
	files = append(files, FileEntry{path.Join(output, "variables.cue"), variablesSrc})
	for _, row := range monitoring.Rows {
		rowDecls := append(
			[]ast.Decl{Package(row.Name), Imports(CuemonPath), ast.NewSel(ast.NewIdent(CuemonName), RowVar)},
			row.Src...,
		)
		rowSrc, err := FormatDecls(rowDecls)
		if err != nil {
			return nil, fmt.Errorf("unable to format row '%v' source: %w", row.Name, err)
		}
		files = append(files, FileEntry{path.Join(output, "rows", fmt.Sprintf("%v.cue", row.Name)), rowSrc})
	}
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
	version := versions[0]
	m := monitoringContext{context: fmt.Sprintf(`%v`, version)}
	panelDecls, err := m.createPanel(panel)
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
