package cuemon

import (
	"list"
)

#playground: bool | *false @tag(playground,type=bool) // #playground: true will remove any potential side-effects from dashboards (alerts notifications, tags, etc)

links: [string]: string
tags: [...string]
variables: [string]: #variable
rows: [...#row]

#variable: {
	#constant: {
		value: string
	}

	#custom: {
		label?:     string
		multi:      *true | bool
		includeAll: *multi | bool
		values: [...string]
		current: [...string]
		if !multi {current: list.MinItems(1) & list.MaxItems(1)}
	}

	#query: {
		label?:     string
		dataSource: string
		query:      string
		multi:      *true | bool
		includeAll: *multi | bool
		if includeAll {allValue?: string}

		current: [...string]
		sort: *"alph-asc-icase" | "alph-desc-icase" | "disabled" | "alph-desc" | "num-asc" | "num-desc" | "alph-asc"
	}

	type: "constant" | "custom" | "query"
	if type == "constant" {#constant}
	if type == "custom" {#custom}
	if type == "query" {#query}
}

#row: {
	title?:    string
	collapsed: *false | bool

	columns: *[12, 12] | [...number] & list.MinItems(1)
	heights: *9 | [...number] | number

	panel: [string]: #panel
	#valid: list.Sum(columns) & <=24
}

#panel: {
	#dataSource: {
		type: "vm" | "stackdriver"
		name: string
		parameters?: [string]: _
	}
	#alert: {
		#threshold: {
			ref:         =~"^[^,]+$"
			aggregation: "avg" | "min" | "max" | "sum"
			period:      "1m" | "5m" | "10m" | "15m" | "30m" | "1h" | "2h" | "3h" | "6h"
			start:       "now" | "now-1m" | "now-5m"
			operator:    "above" | "below" | "in" | "notin"
			if operator == "above" || operator == "below" {parameter: number}
			if operator == "in" || operator == "notin" {parameter: [number, number]}
		}

		name:                string
		message?:            string
		noDataState:         *"no_data" | "keep_state" | "ok" | "alerting"
		executionErrorState: *"keep_state" | "alerting"
		pendingPeriod:       *"5m" | string
		frequency:           *"1m" | string
		thresholds: [...#threshold]
		channels: [...string]
		tags: [string]: string
	}
	#metric: {
		query:  string
		legend: string
		hide:   *false | bool
		parameters?: [string]: _
	}
	#widget: {
		#graph: {
			unit:           string
			yMin:           *null | number
			yMax:           *null | number
			pointSize:      *0 | int & >=0
			lineWidth:      *1 | int & >=0
			fillOpacity:    *0 | int & >=0
			stack:          *false | bool
			nullValue:      *"null" | "null_as_zero" | "connected"
			sort:           *"current" | "avg" | "max" | "min" | "total" | null
			sortDesc:       *false | bool
			legendPosition: *"table_right" | "table_bottom" | "list_right" | "list_bottom" | "none"
			legendValues: [...("last" | "mean" | "min" | "avg" | "max" | "total" | "current" | "lastNotNull" | "sum")]
		}
		#stat: {
			textMode:  *"auto" | "value" | "value_and_name" | "name" | "none"
			graphMode: *"area" | "none"
			reduce:    *"lastNotNull" | "last" | "all"
			thresholds: [...{color: string, value: number}]
		}
		#gauge: {}
		#table: {}

		type: "graph" | "stat" | "gauge" | "table"
		if type == "graph" {#graph}
		if type == "stat" {#stat}
		if type == "gauge" {#gauge}
		if type == "table" {#table}
	}
	#grid: {x?: number, y?: number, width?: number, height?: number}

	grid?:      #grid
	dataSource: #dataSource
	metrics: [...#metric]
	widget: #widget
	alert?: #alert
}
