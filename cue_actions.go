package main

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/tools/trim"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func cueAst(data string) *ast.File {
	cueContext := cuecontext.New()
	return cueContext.CompileString(data).Source().(*ast.File)
}

var packageRegex = regexp.MustCompile(`package \w+`)

func cuePackage(data string) string {
	return strings.SplitN(packageRegex.FindString(data), " ", 2)[1]
}

func buildOverlay(sources map[string]string, target string) (string, []string, map[string]load.Source) {
	hackName := fmt.Sprintf("/overlay-target-%v.cue", rand.Int())
	filenames := make([]string, 0)
	overlay := map[string]load.Source{hackName: load.FromString(target)}
	for filename, source := range sources {
		filenames = append(filenames, filename)
		overlay[filename] = load.FromString(source)
	}
	filenames = append(filenames, hackName)
	return hackName, filenames, overlay
}

func forceTrim(sources map[string]string, target string) (ast.Node, error) {
	cueContext := cuecontext.New()
	targetName, filenames, overlay := buildOverlay(sources, target)
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
		return file, nil
	}
	return nil, fmt.Errorf("unable to trim node against variable")
}

func cuePrettify(decl ast.Decl, flatten bool) ([]ast.Decl, error) {
	embedDecl, ok := decl.(*ast.EmbedDecl)
	if !ok {
		return []ast.Decl{decl}, nil
	}
	structLit, ok := embedDecl.Expr.(*ast.StructLit)
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

func cueConvert(conversions map[string]string, data any, flatten bool) ([]ast.Decl, error) {
	packages := make([]string, 0)
	for _, conversion := range conversions {
		packages = append(packages, cuePackage(conversion))
	}
	uniquePackages := getUnique(packages)
	if len(uniquePackages) > 1 {
		return nil, fmt.Errorf("all conversion sources should belong to the same package: %v", packages)
	}
	pack := "default"
	if len(uniquePackages) > 0 {
		pack = uniquePackages[0]
	}

	input, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json: %w", err)
	}
	target := fmt.Sprintf(`
package %v
Output: (#Conversion & {Input: Data}).Output

Data: %v`, pack, string(input))
	cueContext := cuecontext.New()
	_, filenames, overlay := buildOverlay(conversions, target)
	buildInstances := load.Instances(filenames, &load.Config{Overlay: overlay})
	original, err := cueContext.BuildInstances(buildInstances)
	if err != nil {
		return nil, fmt.Errorf("unable to compile conversion sources and target: %w", err)
	}
	source := original[0].LookupPath(cue.MakePath(cue.Str("Output"))).Eval()
	converted := fmt.Sprintf("%v", source)
	return cuePrettify(cueAst(converted).Decls[0], flatten)
}
