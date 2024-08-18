package lib

import "encoding/json"

type Box struct {
	Id int
	X  int `json:"x"`
	Y  int `json:"y"`
	W  int `json:"w"`
	H  int `json:"h"`
}

type Datasource = JsonRaw[any]

type GrafanaPanel struct {
	Id         int                     `json:"id"`
	Type       string                  `json:"type"`
	Title      string                  `json:"title"`
	Collapsed  bool                    `json:"collapsed"`
	GridPos    Box                     `json:"gridPos"`
	Targets    []json.RawMessage       `json:"targets"`
	Panels     []JsonRaw[GrafanaPanel] `json:"panels"`
	Datasource Datasource              `json:"datasource"`
}

type Templating = JsonRaw[struct {
	Name string `json:"name"`
}]

type Link = JsonRaw[any]

type Grafana = JsonRaw[struct {
	Id            int    `json:"id"`
	Uid           string `json:"uid"`
	Title         string `json:"title"`
	SchemaVersion int    `json:"schemaVersion"`
	Links         []Link `json:"links"`
	Templating    struct {
		List []Templating `json:"list"`
	} `json:"templating"`
	Panels []JsonRaw[GrafanaPanel] `json:"panels"`
}]

type CueRow struct {
	Title     string `json:"title"`
	Collapsed bool   `json:"collapsed"`
	Id        int
	Y         int
	Height    *int   `json:"h"`
	Widths    *[]int `json:"w"`
	Groups    []struct {
		Height *int                    `json:"h"`
		Widths *[]int                  `json:"w"`
		Panels []JsonRaw[GrafanaPanel] `json:"panels"`
	} `json:"groups"`
}
