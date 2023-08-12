package lib

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/tools/trim"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

func CueAst(data string) (*ast.File, error) {
	return parser.ParseFile("", data)
}

func FormatDecls(decl []ast.Decl) string {
	return FormatNode(File(decl))
}

func FormatNode(node ast.Node) string {
	content, err := format.Node(node, format.Simplify())
	if err != nil {
		log.Fatalf("unable to format declarations: %v", err)
	}
	return string(content)
}

func matchLabels(a, b ast.Node) bool {
	identA, okIdentA := a.(*ast.Ident)
	identB, okIdentB := b.(*ast.Ident)
	if okIdentA && okIdentB && identA.Name == identB.Name {
		return true
	}
	stringA, okStringA := a.(*ast.BasicLit)
	stringB, okStringB := b.(*ast.BasicLit)
	if okStringA && okStringB && stringA.Value == stringB.Value {
		return true
	}
	return false
}

func matchPath(node ast.Node, path ...ast.Node) ast.Node {
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

func findField(decls []ast.Decl, path ...ast.Node) (ast.Node, ast.Node) {
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

func CueConvert(variable string, conversions []string, jsonRaw []byte, jsonContext []byte, flatten bool) ([]ast.Decl, error) {
	pack, err := singleCuePackage(conversions)
	if err != nil {
		return nil, fmt.Errorf("invalid packages assignment: %w", err)
	}
	target := fmt.Sprintf(`
package %v
Output: (#Conversion & {Input: Data}).Output

Data: %v
%v`, pack, string(jsonRaw), string(jsonContext))
	cueContext := cuecontext.New()
	_, filenames, overlay := buildOverlay(conversions, target)
	buildInstances := load.Instances(filenames, &load.Config{Overlay: overlay})
	original, err := cueContext.BuildInstances(buildInstances)
	if err != nil {
		return nil, fmt.Errorf("unable to compile conversion sources and target: %w", err)
	}
	source := original[0].LookupPath(cue.MakePath(cue.Str("Output"))).Syntax(cue.Concrete(true))
	converted, err := format.Node(source, format.Simplify())
	if err != nil {
		return nil, fmt.Errorf("unable to get concrete value of conversion: %w", err)
	}
	trimmed, err := ForceTrim(variable, conversions, string(converted))
	if err != nil {
		return nil, fmt.Errorf("unable to trim concrete value: %w", err)
	}
	pretty, err := CuePrettify(trimmed, flatten)
	if err != nil {
		return nil, fmt.Errorf("unable to prettify value: %w", err)
	}
	if _, ok := pretty[0].(*ast.BottomLit); ok {
		return nil, fmt.Errorf("unable to convert value: %v", string(jsonRaw))
	}
	return pretty, nil
}
