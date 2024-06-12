package lib

import (
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/tools/trim"
)

type ReductionRule struct {
	Reduction func(path []string, value ast.Expr) (string, ast.Expr, bool)
}

type ReductionRules []ReductionRule

func (r ReductionRules) reduceField(path []string, decl ast.Decl) (ast.Decl, bool) {
	field, ok := decl.(*ast.Field)
	if !ok {
		return decl, true
	}
	label, ok := labelString(field.Label)
	if !ok {
		return decl, true
	}
	value := field.Value
	pathWithLabel := append(path, label)
	r.reduceAst(pathWithLabel, []ast.Decl{field.Value})
	for _, rule := range r {
		newLabel, newValue, ok := rule.Reduction(pathWithLabel, value)
		if !ok {
			continue
		}
		if newValue == nil {
			return nil, false
		}
		value = newValue
		label = newLabel
	}
	field = FieldIdent(label, value)
	return field, true
}

func (r ReductionRules) reduceAst(path []string, decls []ast.Decl) []ast.Decl {
	result := make([]ast.Decl, 0)
	for _, decl := range decls {
		if structLit, ok := decl.(*ast.StructLit); ok {
			structLit.Elts = r.reduceAst(path, structLit.Elts)
			result = append(result, structLit)
		} else if arrayLit, ok := decl.(*ast.ListLit); ok {
			for _, decl := range arrayLit.Elts {
				fmt.Printf("")
				r.reduceAst(path, []ast.Decl{decl})
			}
			result = append(result, arrayLit)
		} else if field, ok := decl.(*ast.Field); ok {
			field, ok := r.reduceField(path, field)
			if ok {
				result = append(result, field)
			}
		}
	}
	return result
}

func (r ReductionRules) ReduceAst(decls []ast.Decl) []ast.Decl {
	return r.reduceAst(nil, decls)
}

func CueAst(data string) (*ast.File, error) {
	return parser.ParseFile("", data)
}

func FormatDecls(decl []ast.Decl) (string, error) {
	return FormatNode(File(decl))
}

func FormatNode(node ast.Node) (string, error) {
	content, err := format.Node(node, format.Simplify())
	if err != nil {
		return "", fmt.Errorf("unable to format declarations: %w", err)
	}
	return string(content), nil
}

func labelString(a ast.Node) (string, bool) {
	if ident, ok := a.(*ast.Ident); ok {
		return ident.Name, true
	}
	if s, ok := a.(*ast.BasicLit); ok {
		return s.Value[1 : len(s.Value)-1], true
	}
	return "", false
}

func matchLabels(a string, b ast.Node) bool {
	bString, bOk := labelString(b)
	if !bOk {
		return false
	}
	return a == bString
}

func matchPath(node ast.Node, path ...string) ast.Node {
	if _, ok := node.(*ast.Field); ok {
		node = ast.NewStruct(node)
	}
	for _, component := range path {
		s, ok := node.(*ast.StructLit)
		if !ok {
			return nil
		}
		var target ast.Node
		for _, element := range s.Elts {
			field, ok := element.(*ast.Field)
			if !ok {
				continue
			}
			if !matchLabels(component, field.Label) {
				continue
			}
			target = field.Value
		}
		if target == nil {
			return nil
		}
		node = target
	}
	return node
}

func findField(decls []ast.Decl, path ...string) (ast.Node, ast.Node) {
	for _, decl := range decls {
		node := matchPath(decl, path...)
		if node != nil {
			return decl, node
		}
	}
	return nil, nil
}

var packageRegex = regexp.MustCompile(`package \w+`)

func cuePackage(data string) string {
	return strings.SplitN(packageRegex.FindString(data), " ", 2)[1]
}

func singleCuePackage(sources []string) (string, error) {
	packages := make([]string, 0)
	for _, conversion := range sources {
		packages = append(packages, cuePackage(conversion))
	}
	uniquePackages := Unique(packages)
	if len(uniquePackages) == 0 {
		return "defualt", nil
	}
	if len(uniquePackages) > 1 {
		return "", fmt.Errorf("all conversion sources should belong to the same package: %v", packages)
	}
	return uniquePackages[0], nil
}

func buildOverlay(sources []string, target string) (string, []string, map[string]load.Source) {
	hackName := fmt.Sprintf("/overlay-target-%v.cue", rand.Int())
	filenames := make([]string, 0)
	overlay := map[string]load.Source{hackName: load.FromString(target)}
	for _, source := range sources {
		filename := fmt.Sprintf("/overlay-source-%v.cue", rand.Int())
		filenames = append(filenames, filename)
		overlay[filename] = load.FromString(source)
	}
	filenames = append(filenames, hackName)
	return hackName, filenames, overlay
}

func ForceTrim(variable string, sources []string, target string) (ast.Decl, error) {
	pack, err := singleCuePackage(sources)
	if err != nil {
		return nil, fmt.Errorf("invalid packages assignment: %w", err)
	}
	cueContext := cuecontext.New()
	targetName, filenames, overlay := buildOverlay(sources, fmt.Sprintf("package %v\n%v & %v", pack, variable, target))
	originalBuildInstances := load.Instances(filenames, &load.Config{Overlay: overlay})
	original, err := cueContext.BuildInstances(originalBuildInstances)
	if err != nil {
		return nil, fmt.Errorf("unable to build original CUE instances: %w (%v)", err, target)
	}
	err = trim.Files(originalBuildInstances[0].Files, original[0].Value(), &trim.Config{})
	var trimmedHackBytes []byte
	for _, file := range originalBuildInstances[0].Files {
		if file.Filename != targetName {
			continue
		}
		trimmedHackBytes, err = format.Node(file)
		if err != nil {
			return nil, fmt.Errorf("unable to format trimmed hack-file: %w", err)
		}
		break
	}
	if len(trimmedHackBytes) == 0 {
		return nil, fmt.Errorf("unable to get trimmed hack-file")
	}
	targetName, filenames, overlay = buildOverlay(sources, string(trimmedHackBytes))
	trimmedBuildInstances := load.Instances(filenames, &load.Config{Overlay: overlay})
	for _, file := range trimmedBuildInstances[0].Files {
		if file.Filename != targetName {
			continue
		}
		return file.Decls[1].(*ast.EmbedDecl).Expr.(*ast.BinaryExpr).Y, nil
	}
	return nil, fmt.Errorf("unable to trim node against variable")
}

func CuePrettify(decl ast.Decl, flatten bool) ([]ast.Decl, error) {
	embedDecl, ok := decl.(*ast.EmbedDecl)
	if ok {
		decl = embedDecl.Expr
	}
	structLit, ok := decl.(*ast.StructLit)
	if !ok {
		return []ast.Decl{decl}, nil
	}
	var final []ast.Decl
	if !flatten {
		final = structLit.Elts
	} else {
		for _, element := range structLit.Elts {
			switch node := element.(type) {
			case *ast.Field:
				if nested, ok := node.Value.(*ast.StructLit); ok {
					for _, inner := range nested.Elts {
						final = append(final, &ast.Field{Label: node.Label, Value: ast.NewStruct(inner)})
					}
				} else {
					final = append(final, element)
				}
			default:
				final = append(final, element)
			}
		}
	}
	return final, nil
}

func CueFilter(decls []ast.Decl, filterNames []string) []ast.Decl {
	filtered := make([]ast.Decl, 0)
	for _, decl := range decls {
		if field, ok := decl.(*ast.Field); ok {
			label, ok := labelString(field.Label)
			if ok && slices.Contains(filterNames, label) {
				continue
			}
		}
		filtered = append(filtered, decl)
	}
	return filtered
}

func CueConvert(
	conversions []string,
	context map[string]string,
	trimAgainstVar string,
	mappingName string,
	fieldName string,
	flatten bool,
) ([]ast.Decl, error) {
	pack, err := singleCuePackage(conversions)
	if err != nil {
		return nil, fmt.Errorf("invalid packages assignment: %w", err)
	}
	lets := make([]string, 0)
	init := make([]string, 0)
	for key, value := range context {
		lets = append(lets, fmt.Sprintf("let var_%v=%v", key, value))
		init = append(init, fmt.Sprintf("%v: var_%v", key, key))
	}
	template := `
package %v

%v

%v: (%v & {%v}).%v`
	target := fmt.Sprintf(
		template,
		pack,
		strings.Join(lets, "\n"),
		fieldName,
		mappingName,
		strings.Join(init, ", "),
		fieldName,
	)
	cueContext := cuecontext.New()
	_, filenames, overlay := buildOverlay(conversions, target)
	buildInstances := load.Instances(filenames, &load.Config{Overlay: overlay})
	original, err := cueContext.BuildInstances(buildInstances)
	if err != nil {
		fmt.Printf("%v\n", conversions)
		return nil, fmt.Errorf("unable to compile conversion sources and target: %w", err)
	}
	source := original[0].LookupPath(cue.MakePath(cue.Str(fieldName))).Syntax(cue.Concrete(true))
	source = replaceHashFields(source)
	converted, err := format.Node(source, format.Simplify())
	if err != nil {
		return nil, fmt.Errorf("unable to get concrete value of conversion: %w", err)
	}
	trimmed, err := ForceTrim(trimAgainstVar, conversions, string(converted))
	if err != nil {
		return nil, fmt.Errorf("unable to trim concrete value: %w", err)
	}
	pretty, err := CuePrettify(removeEmptyFields(trimmed), flatten)
	if err != nil {
		return nil, fmt.Errorf("unable to prettify value: %w", err)
	}
	if _, ok := pretty[0].(*ast.BottomLit); ok {
		return nil, fmt.Errorf("unable to convert value: %v", context)
	}
	return pretty, nil
}

func isEmpty(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.StructLit:
		return len(n.Elts) == 0
	case *ast.ListLit:
		return len(n.Elts) == 0
	}
	return false
}

func removeEmptyFields(node ast.Decl) ast.Decl {
	if listLit, ok := node.(*ast.ListLit); ok {
		for _, element := range listLit.Elts {
			removeEmptyFields(element)
		}
	} else if structLit, ok := node.(*ast.StructLit); ok {
		replaced := make([]ast.Decl, 0)
		for _, decl := range structLit.Elts {
			if field, ok := decl.(*ast.Field); ok {
				removeEmptyFields(field.Value)
				if isEmpty(field.Value) {
					continue
				}
			}
			replaced = append(replaced, decl)
		}
		structLit.Elts = replaced
		return structLit
	}
	return node

}

func replaceHashFields(node ast.Node) ast.Node {
	if listLit, ok := node.(*ast.ListLit); ok {
		for _, element := range listLit.Elts {
			replaceHashFields(element)
		}
	} else if structLit, ok := node.(*ast.StructLit); ok {
		replaced := make([]ast.Decl, 0)
		for _, decl := range structLit.Elts {
			if field, ok := decl.(*ast.Field); ok {
				replaceHashFields(field.Value)
				label, ok := labelString(field.Label)
				if ok && strings.HasPrefix(label, "#") {
					replaced = append(replaced, FieldIdent(label, field.Value))
					continue
				}
			}
			replaced = append(replaced, decl)
		}
		structLit.Elts = replaced
		return structLit
	}
	return node
}
