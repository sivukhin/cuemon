package lib

import "encoding/json"

type Box struct {
	Id int
	X  int `json:"x"`
	Y  int `json:"y"`
	W  int `json:"w"`
	H  int `json:"h"`
}

type GrafanaPanel struct {
	Type      string                  `json:"type"`
	Title     string                  `json:"title"`
	Collapsed bool                    `json:"collapsed"`
	GridPos   Box                     `json:"gridPos"`
	Targets   []json.RawMessage       `json:"targets"`
	Panels    []JsonRaw[GrafanaPanel] `json:"panels"`
}

type Templating = JsonRaw[struct {
	Name string `json:"name"`
}]

type Grafana = JsonRaw[struct {
	Id            int `json:"id"`
	SchemaVersion int `json:"schemaVersion"`
	Templating    struct {
		List []Templating `json:"list"`
	} `json:"templating"`
	Panels []JsonRaw[GrafanaPanel] `json:"panels"`
}]
