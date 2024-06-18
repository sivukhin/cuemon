package lib

import (
	"encoding/json"
)

type JsonRaw[T any] struct {
	Value T
	Raw   []byte
}

func (d *JsonRaw[T]) UnmarshalJSON(data []byte) error {
	var payload map[string]any
	_ = json.Unmarshal(data, &payload)
	d.Raw, _ = json.Marshal(payload)
	return json.Unmarshal(data, &d.Value)
}
