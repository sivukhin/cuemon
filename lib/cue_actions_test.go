package lib

import (
	"testing"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		trimVar := "{ #meta?: { id: number, uid: number }, if #meta != _|_ { id: #meta.id } }"
		conversions := []string{
			`package main
#Conversion: { Input: { id: number }, Output: { #meta: { id: Input.id, uid: -id }, Input } }
`,
		}
		result, err := CueConvert(
			conversions,
			map[string]string{"Input": `{"id": 123}`},
			trimVar,
			"#Conversion",
			"Output",
			false,
		)
		require.Nil(t, err)
		node, err := format.Node(File([]ast.Decl{ast.NewStruct(AsAny(result)...)}))
		require.Nil(t, err)
		require.Equal(t, `{
	#meta: {
		id:  123
		uid: -123
	}
}
`, string(node))
	})
	t.Run("simple", func(t *testing.T) {
		conversion := []string{
			`package main
#Conversion: { Input: number, Output: Input + Input }
`,
		}
		result, err := CueConvert(
			conversion,
			map[string]string{"Input": "123"},
			"number",
			"#Conversion",
			"Output",
			false,
		)
		require.Nil(t, err)
		node, err := format.Node(File(result))
		require.Nil(t, err)
		require.Equal(t, "246\n", string(node))
	})
	t.Run("complex", func(t *testing.T) {
		conversions := []string{
			`package main

#Def: {
	X: number
	Y: string
	Z: {
		A: number
		B: bool
	}
}
`,
			`package main
#Conversion: { Input: #Def, Output: {
	X: 10 * Input.X
	Y: len(Input.Y)
	A: -Input.Z.A
	if Input.Z.B {
		Win: true
	}
}}
`,
		}
		result, err := CueConvert(
			conversions,
			map[string]string{"Input": `{"X":1,"Y":"hello","Z":{"A":5,"B":true}}`},
			"_",
			"#Conversion",
			"Output",
			true,
		)
		require.Nil(t, err)
		node, err := format.Node(File(result), format.Simplify())
		require.Nil(t, err)
		require.Equal(t, `X:   10
Y:   5
Win: true
A:   -5
`, string(node))
	})
}

func TestPrettify(t *testing.T) {
	ast, err := CueAst(`{
	A: {
		C: "123"
		D: "456"
	}
	B: 2
}`)
	require.Nil(t, err)
	declarations, err := CuePrettify(ast.Decls[0], true)
	require.Nil(t, err)
	pretty, err := format.Node(File(declarations), format.Simplify())
	require.Nil(t, err)
	require.Equal(t, `A: C: "123"
A: D: "456"
B: 2
`, string(pretty))
}

func TestForceTrim(t *testing.T) {
	result, err := ForceTrim("#D", []string{
		`package main
#D: {
	type: "A" | "B"
	if type == "A" {
		value: number | *0
	}
}`,
	}, `{type: "A", value: 0}`)
	require.Nil(t, err)
	node, err := format.Node(result)
	t.Log(string(node))
	require.Nil(t, err)
	require.Equal(t, `{type: "A"}`, string(node))
}
