package main

import (
	"cuelang.org/go/cue/format"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		result, err := cueConvert(map[string]string{
			"/a.cue": `package main
#Conversion: { Input: number, Output: Input + Input }
`,
		}, 123, false)
		require.Nil(t, err)
		node, err := format.Node(File(result))
		require.Nil(t, err)
		require.Equal(t, "246\n", string(node))
	})
	t.Run("complex", func(t *testing.T) {
		result, err := cueConvert(map[string]string{
			"/def.cue": `package main

#Def: {
	X: number
	Y: string
	Z: {
		A: number
		B: bool
	}
}
`,
			"/a.cue": `package main
#Conversion: { Input: #Def, Output: {
	X: 10 * Input.X
	Y: len(Input.Y)
	A: -Input.Z.A
	if Input.Z.B {
		Win: true
	}
}}
`,
		}, map[string]any{"X": 1, "Y": "hello", "Z": map[string]any{"A": 5, "B": true}}, false)
		require.Nil(t, err)
		node, err := format.Node(File(result))
		require.Nil(t, err)
		require.Equal(t, `X:   10
Y:   5
Win: true
A:   -5
`, string(node))
	})
}

func TestPrettify(t *testing.T) {
	declarations, err := cuePrettify(cueAst(`{
	A: {
		C: "123"
		D: "456"
	}
	B: 2
}`).Decls[0], true)
	require.Nil(t, err)
	pretty, err := format.Node(File(declarations), format.Simplify())
	require.Nil(t, err)
	require.Equal(t, `A: C: "123"
A: D: "456"
B: 2
`, string(pretty))
}

func TestForceTrim(t *testing.T) {
	result, err := forceTrim(map[string]string{
		"/a.cue": `package main
#D: {
	type: "A" | "B"
	if type == "A" {
		value: number | *0
	}
}`,
	}, `package main 

#D & {type: "A", value: 0}
`)
	require.Nil(t, err)
	node, err := format.Node(result)
	t.Log(string(node))
	require.Nil(t, err)
	require.Equal(t, `package main

#D & {type: "A"}
`, string(node))
}
