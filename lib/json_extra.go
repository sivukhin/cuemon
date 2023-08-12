package lib

import (
	"encoding/json"
)

type JsonRaw[T any] struct {
	Value T
	Raw   []byte
}

func (d *JsonRaw[T]) UnmarshalJSON(data []byte) error {
	d.Raw = data
	return json.Unmarshal(data, &d.Value)
}
