package cuemon

import (
	"list"
	"strings"
	"strconv"
)

#monRow: {
	type:      "row"
	title?:    string
	collapsed: bool | *false
	if #grafanaVersion == "v7" {datasource: _ | *null}

	rowH=h: number
	rowW=w: [...number]

	groups: [...{
		h: number | *rowH
		w: [...number] | *rowW

		panels: [...#monPanel]
	}]
}

#monPanel: #panel & {#panel}
#monPanel: {
	#text?: {
		title:       string | *""
		html?:       string
		markdown?:   string
		transparent: bool | *false
	}
	if #text != _|_ {
		type:          "text"
		title:         #text.title
		transparent:   #text.transparent
		pluginVersion: #pluginVersion[type]
		if #text.html != _|_ {
			options: content: #text.html
			options: mode:    "html"
		}
		if #text.markdown != _|_ {
			options: content: #text.markdown
			options: mode:    "markdown"
		}
		if #grafanaVersion == "v10" {
			options: code: language:        _ | *"plaintext"
			options: code: showLineNumbers: _ | *false
			options: code: showMiniMap:     _ | *false
			error:    _ | *false
			editable: _ | *true
		}
		if #grafanaVersion == "v7" {
			datasource: _ | *null
		}
	}
}

#monPanel: {
	// set threshold ranges for colors of the text in stat plugin
	#thresholds?: {mode: "percentage" | *"absolute", steps: [...{color: string, value: number | *null}]}

	if #thresholds != _|_ {
		fieldConfig: defaults: thresholds: mode:  #thresholds.mode
		fieldConfig: defaults: thresholds: steps: #thresholds.steps
		if len(#thresholds.steps) > 1 {
			fieldConfig: defaults: color: mode: _ | *"thresholds"
		}
	}
}
#monPanel: {
	// set units of the graph
	#unit?: string

	if #unit != _|_ {
		fieldConfig: defaults: unit: _ | *#unit
	}
}

#monPanel: {
	// set min/max value for the y-axis of the graph
	#rangeY?: {min?: number, max?: number}
}

#monPanel: {
	#stat?: {
		title: string
		reducer: [...("lastNotNull" | "last" | "min" | "max" | "avg")] | *["lastNotNull"]
	}

	if #stat != _|_ {
		type:  "stat"
		title: #stat.title
		options: reduceOptions: calcs: _ | *#stat.reducer
	}
}

#monPanel: {
	#graph?: {
		#values: "avg" | "min" | "max" | "current" | "total"

		legend: type: *"table" | "list"
		legend: values: [...#values] | *["current"]
		legend: placement: "bottom" | *"right"
		legend: sortBy:    #values | *legend.values[0]
		legend: sortHow:   "asc" | *"desc"
	}

	if #graph != _|_ {
		type: "graph"
		xaxis: mode:         _ | *"time"
		xaxis: show:         _ | *true
		tooltip: shared:     _ | *true
		tooltip: sort:       _ | *2
		tooltip: value_type: _ | *"individual"
		bars:          _ | *false
		dashLength:    _ | *10
		dashes:        _ | *false
		fill:          _ | *0
		fillGradient:  _ | *0
		hiddenSeries:  _ | *false
		lines:         _ | *true
		linewidth:     _ | *1
		nullPointMode: _ | *"null"
		percentage:    _ | *false
		pluginVersion: _ | *#pluginVersion[type]
		pointradius:   _ | *2
		points:        _ | *false
		renderer:      _ | *"flot"
		spaceLength:   _ | *10
		stack:         _ | *false

		legend: avg:     _ | *list.Contains(#graph.legend.values, "avg")
		legend: max:     _ | *list.Contains(#graph.legend.values, "max")
		legend: min:     _ | *list.Contains(#graph.legend.values, "min")
		legend: current: _ | *list.Contains(#graph.legend.values, "current")
		legend: total:   _ | *list.Contains(#graph.legend.values, "total")

		legend: rightSide: _ | *(#graph.legend.placement == "right")
		legend: sort:      _ | *#graph.legend.sortBy
		legend: sortDesc:  _ | *(#graph.legend.sortHow == "desc")

		legend: alignAsTable: _ | *(#graph.legend.type == "table")
		legend: values:       _ | *true
		legend: show:         _ | *true
	}
}

#monPanel: {
	#alphabet: [for x in list.Range(10, 256+10, 1) {strings.ToUpper(strconv.FormatInt(x, 36))}]

	type:          string
	#d=datasource: #datasource

	targets: [for i, target in targets {
		refId:  _ | *#alphabet[i]
		#prom?: string
		if #prom != _|_ {
			expr:     _ | *#prom
			interval: _ | *""
		}
		#mqlScript?: {
			query:              string
			aliasBy:            string
			projectName:        string
			alignmentPeriod:    string | *"cloud-monitoring-auto"
			crossSeriesReducer: *"REDUCE_MEAN" | "REDUCE_SUM"
			perSeriesAligner:   *"ALIGN_MEAN" | "ALIGN_INTERPOLATE" | "ALIGN_NEXT_OLDER" | "ALIGN_RATE"
		}
		if #mqlScript != _|_ {
			metricQuery: aliasBy:            _ | *#mqlScript.aliasBy
			metricQuery: alignmentPeriod:    _ | *#mqlScript.alignmentPeriod
			metricQuery: crossSeriesReducer: _ | *#mqlScript.crossSeriesReducer
			metricQuery: perSeriesAligner:   _ | *#mqlScript.perSeriesAligner
			metricQuery: projectName:        _ | *#mqlScript.projectName
			metricQuery: query:              _ | *#mqlScript.query

			metricQuery: editorMode: _ | *"mql"
			metricQuery: metricKind: _ | *""
			metricQuery: metricType: _ | *""
			metricQuery: unit:       _ | *""
			metricQuery: valueType:  _ | *""

			queryType: _ | *"metrics"
		}

		#format: string | *""
		#hide:   bool | *false

		datasource?:  _ | *#d
		legendFormat: _ | *#format
		hide:         _ | *#hide
		if type == "stat" || type == "table" {
			exemplar: _ | *false
			instant:  _ | *true
		}
		if type == "timeseries" || type == "graph" {
			instant: _ | *false
			//			range:    _ | *true
			exemplar: _ | *true
		}
	}]
}

#monPanel: {
	panelH=h?: number
	if panelH != _|_ {gridPos: h: _ | *panelH}

	type:          string
	pluginVersion: _ | *#pluginVersion[type]

	if type == "stat" && #grafanaVersion == "v7" {
		options: alertThreshold: _ | *true
		options: reduceOptions: fields: _ | *""
		options: reduceOptions: values: _ | *false

		options: colorMode:   _ | *"value"
		options: graphMode:   _ | *"area"
		options: justifyMode: _ | *"auto"
		options: orientation: _ | *"auto"
		options: textMode:    _ | *"auto"

		description: _ | *""
		timeShift:   _ | *null
		timeFrom:    _ | *null
		steppedLine: _ | *false
		yaxis: align: _ | *false
		xaxis: mode:  _ | *"time"
		xaxis: show:  _ | *true
		nullPointMode: _ | *"null"
		percentage:    _ | *false
		points:        _ | *false
		pointradius:   _ | *0
		linewidth:     _ | *1
		lines:         _ | *true
		bars:          _ | *false
		dashes:        _ | *false
		dashLength:    _ | *10
		fill:          _ | *0
		fillGradient:  _ | *0
		hiddenSeries:  _ | *false
		legend: alignAsTable: _ | *false
		legend: avg:          _ | *false
		legend: current:      _ | *false
		legend: max:          _ | *false
		legend: min:          _ | *false
		legend: rightSide:    _ | *false
		legend: show:         _ | *false
		legend: sortDesc:     _ | *false
		legend: total:        _ | *false
		legend: values:       _ | *false
		renderer:    _ | *"flot"
		spaceLength: _ | *10
		stack:       _ | *false
		tooltip: shared:     _ | *true
		tooltip: sort:       _ | *2
		tooltip: value_type: _ | *"individual"
	}
	if type == "graph" && #grafanaVersion == "v7" {
		steppedLine: _ | *false
		yaxis: align: _ | *false
	}
}
