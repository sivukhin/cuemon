package lib

import (
	"strconv"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/token"
)

func Int(value int) *ast.BasicLit {
	return ast.NewLit(token.INT, strconv.Itoa(value))
}

func Package(name string) *ast.Package {
	return &ast.Package{Name: ast.NewIdent(name)}
}

func Imports(paths ...string) *ast.ImportDecl {
	imports := make([]*ast.ImportSpec, 0)
	for _, path := range paths {
		imports = append(imports, ast.NewImport(nil, path))
	}
	return &ast.ImportDecl{Specs: imports}
}

func File(decl []ast.Decl) *ast.File {
	return &ast.File{Decls: decl}
}

func FieldIdent(ident string, value ast.Expr, attrs ...*ast.Attribute) *ast.Field {
	if strings.HasPrefix(ident, "__") {
		return &ast.Field{Label: ast.NewString(ident), Value: value, Attrs: attrs}
	}
	return &ast.Field{Label: ast.NewIdent(ident), Value: value, Attrs: attrs}
}

func IntList(ints []int) *ast.ListLit {
	list := make([]ast.Expr, 0, len(ints))
	for _, i := range ints {
		list = append(list, ast.NewLit(token.INT, strconv.Itoa(i)))
	}
	return ast.NewList(list...)
}

func IdentList(idents []string) *ast.ListLit {
	list := make([]ast.Expr, 0, len(idents))
	for _, ident := range idents {
		list = append(list, ast.NewIdent(ident))
	}
	return ast.NewList(list...)
}

func LineBreak() *ast.CommentGroup {
	return &ast.CommentGroup{List: []*ast.Comment{{Text: ""}}}
}
