package cuemon

import (
	"list"
)

Title: string
Links: [string]: Url: string
Variables: [string]: #Variable
Panels: [string]:    #Panel
Rows: [...#Row]
Tags: [...string]

#Target: {
	Expr:    string
	Legend?: string
	Hide: bool | *false
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
	}
}

#Threshold: { Color: string, Value: number | *null }

#LegendValues: ["min", "avg", "max", "total", "current"]
#Panel: {
	Type:       *"graph" | "stat" | "gauge" | "table"
	Unit:       string | *""
	DataSource: string
	Metrics: [...#Target]
	if Type == "graph" {
		Points: number | *0
		Lines: number | *1
		NullValue: *"null" | "null_as_zero" | "connected"
		Legend: "table_right" | *"table_bottom" | "list_right" | "list_bottom" | "none"
		Values: [...or(#LegendValues)] | *["current"]
	}
	if Type == "stat" {
		TextMode: *"auto" | "value" | "value_and_name" | "name" | "none"
		GraphMode: *"area" | "none"
		Reduce: *"lastNotNull" | "all"
		Thresholds?: [...#Threshold]
	}
	Alert?: {
		Name: string
		Message?: string
		NoDataState: *"no_data" | "keep_state" | "ok" | "alerting"
		ExecutionErrorState: *"keep_state" | "alerting"
		PendingPeriod: string | *"5m"
		Frequency: string | *"1m"
		Notifications: [...=~#"(avg|min|max|sum)\(\w+,(1m|5m|10m|15m),(now|now-1m|now-5m)\) (>|<) (.*)"#]
		Tags: [string]: string
		Channels: [...string]
	}
}

#Grid: {X?: number, Y?: number, Width?: number, Height?: number}

#Row: {
	Title?:    string
	Columns:   [number, ...number] | *[]
	Heights:   [...number] | number | *9
	Collapsed: bool | *false
	Panel: [string]:     #Panel
	PanelGrid: [string]: #Grid
	#Width: list.Sum(Columns) & 24
}

#Variable: {
	Type: "constant" | "custom" | "query"
	if Type == "constant" {
		Value: string
	}
	if Type == "custom" {
		Values: [...string]
		Multi:      bool | *true
		IncludeAll: bool | *Multi
		Current: [...string]
	}
	if Type == "query" {
		DataSource: string
		Query:      string
		Multi:      bool | *true
		IncludeAll: bool | *Multi
		Current: [...string]
		Sort: (#GrafanaTemplate & {type: "query"}).sort
	}
}
