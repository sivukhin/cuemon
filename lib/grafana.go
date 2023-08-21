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
	Type            string                   `json:"type"`
	Title           string                   `json:"title"`
	Collapsed       bool                     `json:"collapsed"`
	GridPos         Box                      `json:"gridPos"`
	Targets         []json.RawMessage        `json:"targets"`
	SeriesOverrides []GrafanaSeriesOverrides `json:"seriesOverrides"`
	Panels          []JsonRaw[GrafanaPanel]  `json:"panels"`
}

type GrafanaSeriesOverrides struct {
	Alias     string  `json:"alias"`
	Dashes    *bool   `json:"dashes"`
	Hidden    *bool   `json:"hiddenSeries"`
	Color     *string `json:"color"`
	Fill      *int    `json:"fill"`
	YAxis     *int    `json:"yaxis"`
	ZIndex    *int    `json:"zindex"`
	LineWidth *int    `json:"linewidth"`
}

type Templating = JsonRaw[struct {
	Name string `json:"name"`
}]

type Grafana = JsonRaw[struct {
	Id            int    `json:"id"`
	Uid           string `json:"uid"`
	SchemaVersion int    `json:"schemaVersion"`
	Templating    struct {
		List []Templating `json:"list"`
	} `json:"templating"`
	Panels []JsonRaw[GrafanaPanel] `json:"panels"`
}]
