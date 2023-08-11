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
	Expr: string
	Legend: string
	StackDriver?: {
		Reducer: string
		Filters: [...string]
		GroupBy: [...string]
		Aligner: string
		Project: string
		Uint: string | *""
		Value: string | *""
	}
}

#Panel: {
	Type: *"graph" | "stat" | "table"
	Unit?:       string
	DataSource: string
	Metrics: [...#Target]
}

#Grid: {X?: number, Y?: number, Width?: number, Height?: number}

#Row: {
	Title?: string
	Columns: [number, ...number] | *[]
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
