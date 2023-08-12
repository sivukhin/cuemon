package lib

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonExtra(t *testing.T) {
	type B = JsonRaw[struct {
		Value int `json:"value"`
	}]
	type A = JsonRaw[struct {
		A     string `json:"a"`
		Items []B    `json:"items"`
	}]
	var a A
	jsonBytes := []byte(`{"a":"123","items":[{"value":1},{},{"value":3}]}`)
	err := json.Unmarshal(jsonBytes, &a)
	require.Nil(t, err)
	require.Equal(t, jsonBytes, a.Raw)
	require.Equal(t, a.Value.A, "123")
	require.Len(t, a.Value.Items, 3)
	require.Equal(t, a.Value.Items[0].Value.Value, 1)
	require.Equal(t, a.Value.Items[0].Raw, []byte(`{"value":1}`))
	require.Equal(t, a.Value.Items[1].Value.Value, 0)
	require.Equal(t, a.Value.Items[1].Raw, []byte(`{}`))
	require.Equal(t, a.Value.Items[2].Value.Value, 3)
	require.Equal(t, a.Value.Items[2].Raw, []byte(`{"value":3}`))
}
