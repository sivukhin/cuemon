package main

import (
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/pkg/strings"
	"cuemon/src"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
)

func asAny[T any](array []T) []any {
	a := make([]any, 0, len(array))
	for _, x := range array {
		a = append(a, x)
	}
	return a
}

var (
	//go:embed grafana.cue
	Grafana string
	//go:embed mon.cue
	Mon string
	//go:embed mon2grafana.cue
	Mon2Grafana string
	//go:embed grafana2mon/meta.cue
	Meta2Monitoring string
	//go:embed grafana2mon/template.cue
	Template2Monitoring string
	//go:embed grafana2mon/row.cue
	Row2Monitoring string
	//go:embed grafana2mon/panel.cue
	Panel2Monitoring string
	//go:embed grafana2mon/target.cue
	Target2Monitoring string
)

func has(data any, path ...string) (any, bool) {
	for _, element := range path {
		next, ok := data.(map[string]any)[element]
		if !ok {
			return nil, false
		}
		data = next
	}
	return data, true
}

func hasArray(data any, path ...string) ([]any, bool) {
	data, ok := has(data, path...)
	if !ok {
		return nil, false
	}
	array, ok := data.([]any)
	return array, ok
}

func hasString(data any, path ...string) (string, bool) {
	data, ok := has(data, path...)
	if !ok {
		return "", false
	}
	str, ok := data.(string)
	return str, ok
}

func getBool(data any, path ...string) bool {
	data, ok := has(data, path...)
	if !ok {
		return false
	}
	value, ok := data.(bool)
	return value && ok
}

func getString(data any, path ...string) string {
	str, _ := hasString(data, path...)
	return str
}

var matchIgnoredSymbols = regexp.MustCompile("[^a-zA-Z0-9]")
var matchUnderscores = regexp.MustCompile("_+")

func toFilename(s string) string {
	s = matchIgnoredSymbols.ReplaceAllString(s, "_")
	s = matchUnderscores.ReplaceAllString(s, "_")
	return strings.ToLower(s)
}

func createRow(pack string, row any, panels []any) ([]ast.Decl, error) {
	decls := []ast.Decl{
		src.Package(pack),
		src.Imports("github.com/sivukhin/cuemon"),
		ast.NewSel(ast.NewIdent("cuemon"), "#Row"),
		src.LineBreak(),
	}
	meta, err := src.CueConvert("#Row", []string{Grafana, Mon, Row2Monitoring}, row, true)
	if err != nil {
		return nil, fmt.Errorf("unable to convert row meta: %w", err)
	}
	decls = append(decls, meta...)

	grid := make([]src.Grid, 0)
	for _, panel := range panels {
		panelJson, err := json.Marshal(panel)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal single panel: %w", err)
		}
		var container struct {
			Grid src.Grid `json:"gridPos"`
		}
		err = json.Unmarshal(panelJson, &container)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal single panel: %w", err)
		}
		grid = append(grid, container.Grid)
	}
	layout := src.AnalyzeGrid(grid)

	decls = append(decls, src.FieldIdent("Columns", src.IntList(layout.Columns)))
	if len(layout.Heights) == 1 {
		decls = append(decls, src.FieldIdent("Heights", src.Int(layout.Heights[0])))
	} else {
		decls = append(decls, src.FieldIdent("Heights", src.IntList(layout.Heights)))
	}
	decls = append(decls, src.LineBreak())

	for _, id := range layout.Order {
		panel := panels[id]
		panelDecls, err := src.CueConvert("#Panel", []string{Grafana, Mon, Panel2Monitoring}, panel, false)
		if err != nil {
			return nil, fmt.Errorf("unable to convert panel %v: %w", panel, err)
		}
		title := getString(panel, "title")
		if override, ok := layout.Overrides[id]; ok {
			decls = append(decls, src.FieldIdent("PanelGrid", ast.NewStruct(ast.NewIdent(title), ast.NewStruct(
				src.FieldIdent("Width", src.Int(override.Width)),
				src.FieldIdent("Height", src.Int(override.Height)),
			))))
		}
		declarations := asAny(panelDecls)
		targets, ok := hasArray(panel, "targets")
		if ok {
			targetsDecls := make([]ast.Expr, 0)
			for _, target := range targets {
				targetDecls, err := src.CueConvert("#Target", []string{Grafana, Mon, Target2Monitoring}, target, false)
				if err != nil {
					return nil, fmt.Errorf("unable to convert target from panel %v: %w", panel, err)
				}
				targetsDecls = append(targetsDecls, ast.NewStruct(asAny(targetDecls)...))
			}
			declarations = append(declarations, src.FieldIdent("Metrics", ast.NewList(targetsDecls...)))
		}
		decls = append(decls, src.FieldIdent("Panel", ast.NewStruct(
			ast.NewIdent(title), ast.NewStruct(declarations...),
		)))
		decls = append(decls, src.LineBreak())
	}
	return decls, nil
}

func main() {
	input := flag.String("input", "", "")
	module := flag.String("module", "", "")
	output := flag.String("output", "", "")
	flag.Parse()

	modulePath := strings.Split(*module, "/")
	pack := modulePath[len(modulePath)-1]

	content, err := os.ReadFile(*input)
	if err != nil {
		log.Fatalf("unable to read input file '%v'", *input)
	}
	var data map[string]any
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("unable to unmarshal input file '%v': %v", *input, err)
	}
	meta, err := src.CueConvert("_", []string{Grafana, Meta2Monitoring}, data, true)
	if err != nil {
		log.Fatalf("unable to convert monitoring meta: %v", err)
	}
	variables := []ast.Decl{&ast.Package{Name: ast.NewIdent("variables")}}
	if grafanaTemplating, ok := hasArray(data, "templating", "list"); ok {
		for _, template := range grafanaTemplating {
			name, ok := hasString(template, "name")
			_ = name
			if !ok {
				log.Printf("unable to get variable name '%v'", template)
				continue
			}
			variable, err := src.CueConvert("_", []string{Grafana, Mon, Template2Monitoring}, template, false)
			if err != nil {
				log.Printf("unable to convert variable '%v': %v", template, err)
			}
			declarations := make([]any, 0)
			for _, d := range variable {
				declarations = append(declarations, d)
			}
			variables = append(variables, &ast.Field{Label: ast.NewIdent(name), Value: ast.NewStruct(declarations...)})
		}
	}
	panels := make(map[string][]ast.Decl, 0)
	rows := make([]string, 0)
	_ = panels
	if grafanaPanels, ok := hasArray(data, "panels"); ok {
		name, collapsed := "general", false
		i := 0
		for i < len(grafanaPanels) {
			current := grafanaPanels[i]
			panelType, ok := hasString(current, "type")
			if !ok {
				log.Printf("unable to get panel type at position %v", i)
				continue
			}
			if panelType == "row" {
				title, ok := hasString(current, "title")
				if !ok {
					log.Printf("unable to get panel row title at position %v", i)
					continue
				}
				name = toFilename(title)
				collapsed = getBool(current, "collapsed")
				i += 1
			}
			var rowPanels []any
			if collapsed {
				rowPanels, ok = hasArray(current, "panels")
				if !ok {
					log.Printf("unable to get panels of collapsed row at position %v", i)
					continue
				}
			} else {
				startI := i
				for i < len(grafanaPanels) && getString(grafanaPanels[i], "type") != "row" {
					i += 1
				}
				rowPanels = grafanaPanels[startI:i]
			}
			declarations, err := createRow(name, current, rowPanels)
			if err != nil {
				log.Printf("unable to create row: %v", err)
				continue
			}
			rows = append(rows, name)
			panels[name] = declarations
		}
	}

	imports := []string{"github.com/sivukhin/cuemon"}
	if len(variables) > 0 {
		imports = append(imports, fmt.Sprintf("%v:variables", *module))
		meta = append(meta, src.FieldIdent("Variables", ast.NewIdent("variables")))
	}
	rowList := make([]ast.Expr, 0)
	for _, row := range rows {
		imports = append(imports, fmt.Sprintf("%v/rows:%v", *module, row))
		rowList = append(rowList, ast.NewIdent(row))
	}
	if len(rowList) > 0 {
		meta = append(meta, src.FieldIdent("Rows", ast.NewList(rowList...)))
	}
	decls := append([]ast.Decl{
		src.Package(pack),
		src.Imports(imports...),
		ast.NewSel(ast.NewIdent("cuemon")),
		src.LineBreak(),
	}, meta...)
	type entry struct {
		path    string
		content []byte
	}
	files := make([]entry, 0)
	files = append(files, entry{path.Join(*output, "cue.mod", "module.cue"), []byte(fmt.Sprintf("module: \"%v\"\n", *module))})
	files = append(files, entry{path.Join(*output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "grafana.cue"), []byte(Grafana)})
	files = append(files, entry{path.Join(*output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "mon.cue"), []byte(Mon)})
	files = append(files, entry{path.Join(*output, "cue.mod", "pkg", "github.com", "sivukhin", "cuemon", "mon2grafana.cue"), []byte(Mon2Grafana)})
	dashboardCue, err := format.Node(&ast.File{Decls: decls}, format.Simplify())
	if err != nil {
		log.Printf("unable to format 'dashboard.cue' file: %v", err)
	}
	files = append(files, entry{path.Join(*output, "dashboard.cue"), dashboardCue})
	if len(variables) > 0 {
		variablesCue, err := format.Node(&ast.File{Decls: variables}, format.Simplify())
		if err != nil {
			log.Printf("unable to format 'variables.cue' file: %v", err)
		}
		files = append(files, entry{path.Join(*output, "variables.cue"), variablesCue})
	}
	for _, row := range rows {
		filename := fmt.Sprintf("%v.cue", row)
		rowCue, err := format.Node(&ast.File{Decls: panels[row]}, format.Simplify())
		if err != nil {
			log.Printf("unable to format 'rows/%v': %v", row, err)
		}
		files = append(files, entry{path.Join(*output, "rows", filename), rowCue})
	}

	log.Printf("read to write %v cuemon assets:", len(files))
	for _, file := range files {
		log.Printf("  file %v (%v bytes)", file.path, len(file.content))
	}

	for _, file := range files {
		err = os.MkdirAll(path.Dir(file.path), os.ModePerm)
		if err != nil {
			log.Printf("unable to create directories for file '%v': %v", file.path, err)
			continue
		}
		f, err := os.Create(file.path)
		if err != nil {
			log.Printf("unable to create file '%v': %v", file.path, err)
			continue
		}
		content := file.content
		for len(content) > 0 {
			n, err := f.Write(file.content)
			if err != nil {
				log.Printf("unable to write to file '%v': %v", file.path, err)
				break
			}
			content = content[n:]
		}
		_ = f.Close()
	}
}
