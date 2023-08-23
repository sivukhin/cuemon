package cuemon

import (
	"list"
)

Links: [string]: Url: string
Variables: [string]: #Variable
Panels: [string]:    #Panel
Rows: [...#Row]
Tags: [...string]

#Target: {
	Expr:    string | *""
	Legend?: string
	Hide: bool | *false
	Overrides?: {
		Alias?: string
		Dashes?: bool
		Hidden?: bool
		Fill?: number
		YAxis?: number
		ZIndex?: number
		LineWidth?: number
		Color?: string
	}
	StackDriver?: {
		Reducer: string | *""
		Filters: [...string]
		GroupBy: [...string]
		Aligner:         string | *""
		Project:         string | *""
		AlignmentPeriod: string | *""
		EditorMode: string | *""
		MetricKind:      string | *""
		MetricType:      string | *""
		Unit:            string | *""
		Value:           string | *""
		Preprocessor?:           string
	}
}

#Threshold: { Color: string, Value: number | *null }

#LegendValues: ["last", "mean", "min", "avg", "max", "total", "current", "lastNotNull", "sum"]
#Panel: {
	Type:       "graph" | "stat" | "gauge" | "table" | "timeseries"
	Unit:       string | *""
	DataSource: string
	Metrics: [...#Target]
	if Type == "graph" || Type == "timeseries" {
		YMin: number | *null
		Points: number | *0
		Lines: number | *1
		NullValue: *"null" | "null_as_zero" | "connected"
		Legend: "table_right" | *"table_bottom" | "list_right" | "list_bottom" | "none"
		Values: [...or(#LegendValues)] | *[]
		Sort:           "current" | "avg" | "max" | "min" | "total" | *null
		SortDesc:       bool | *false
	}
	if Type == "stat" {
		TextMode: *"auto" | "value" | "value_and_name" | "name" | "none"
		GraphMode: *"area" | "none"
		Reduce: *"lastNotNull" | "last" | "all"
		Thresholds?: [...#Threshold]
	}
	Alert?: {
		Name: string
		Message?: string
		NoDataState: *"no_data" | "keep_state" | "ok" | "alerting"
		ExecutionErrorState: *"keep_state" | "alerting"
		PendingPeriod: string | *"5m"
		Frequency: string | *"1m"
		Notifications: [...=~#"(avg|min|max|sum)\([^,]+,(1m|5m|10m|15m|1h),(now|now-1m|now-5m)\) (>|<) (.*)"#]
		Tags: [string]: string
		Channels: [...string]
	}
}

#Grid: {X?: number, Y?: number, Width?: number, Height?: number}

#Row: {
	Title?:    string
	Columns:   [number, ...number]
	Heights:   [...number] | number | *9
	Collapsed: bool | *false
	Panel: [string]:     #Panel
	PanelGrid: [string]: #Grid
	_width: list.Sum(Columns) & <=24
}

#Variable: {
	Type: "constant" | "custom" | "query"
	if Type == "constant" {
		Value: string
	}
	if Type == "custom" {
		Label?: string
		Values: [...string]
		Multi:      bool | *true
		IncludeAll: bool | *Multi
		Current: {
			if Multi { [...string] }
			if !Multi { string }
		}
	}
	if Type == "query" {
		Label?: string
		DataSource: string
		Query:      string
		Multi:      bool | *true
		IncludeAll: bool | *Multi
		Current: [...string]
		Sort: (#GrafanaTemplate & {type: "query"}).sort
	}
}
