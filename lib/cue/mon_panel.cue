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
	type: _

	if #thresholds != _|_ && (type != "graph" || #grafanaVersion != "v10") {
		fieldConfig: defaults: thresholds: mode:  _ | *#thresholds.mode
		fieldConfig: defaults: thresholds: steps: _ | *#thresholds.steps
		if len(#thresholds.steps) > 1 {
			fieldConfig: defaults: color: mode: _ | *"thresholds"
		}
	}
	if #thresholds == _|_ && #grafanaVersion == "v10" && type != "graph" {
		fieldConfig: defaults: thresholds: mode:                _ | *"absolute"
		fieldConfig: defaults: thresholds: steps: _ | *[{color: "green"}, {color: "red", value: 80}]
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
	type: _

	if #rangeY != _|_ && type == "timeseries" {
		if #rangeY.min != _|_ {fieldConfig: defaults: min: #rangeY.min}
		if #rangeY.max != _|_ {fieldConfig: defaults: max: #rangeY.max}
	}
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
	#override: {
		alias:         string
		yaxis?:        number
		color:         string
		hashKey?:      string
		legend?:       bool
		hidden?:       bool
		dashes?:       bool
		fillGradient?: number
	}
	#overrides?: [...#override]

	if #overrides != _|_ {
		seriesOverrides: [for i, override in seriesOverrides {
			alias: _ | *#overrides[i].alias
			if #overrides[i].hashKey != _|_ {$$hashKey: _ | *#overrides[i].hashKey}
			if #overrides[i].color != _|_ {color: _ | *#overrides[i].color}
			if #overrides[i].yaxis != _|_ {yaxis: _ | *#overrides[i].yaxis}
			if #overrides[i].hidden != _|_ {hiddenSeries: _ | *#overrides[i].hidden}
			if #overrides[i].dashes != _|_ {dashes: _ | *#overrides[i].dashes}

			// not sure if fillGradient exists for legacy "graph" panels in Grafana v10
			if #overrides[i].fillGradient != _|_ && #grafanaVersion == "v7" {fillGradient: _ | *#overrides[i].fillGradient}
			if #overrides[i].legend != _|_ && #grafanaVersion == "v7" {legend: _ | *#overrides[i].legend}
		}]
	}
}

#monPanel: {
	#defaultGraphPlugin: bool | *true
	#graph?: {
		#values: "min" | "max" | "mean" | "lastNotNull" | "last" | "sum"

		legend: type: *"table" | "list"
		legend: values: [...#values] | *["lastNotNull"]
		legend: showValues: bool | *true
		legend: placement:  "bottom" | *"right"
		legend: sortBy: #values | [...#values] | *"lastNotNull"
		legend: sortHow: "asc" | *"desc"

		leftY: unit:      string | *"short"
		leftY: hashKey?:  string
		rightY: unit:     string | *"short"
		rightY: hashKey?: string

		display: lines: {use: bool | *true, size: number | *1}
		display: points: {use: bool | *false, size: number | *0}
		display: bars: use: bool | *false
		display: nulls: "connected" | *"null"
		display: fill:  number | *0
		display: stack: bool | *false
		tooltip: sort:  number | *2
	}

	if (#graph != _|_ && #grafanaVersion == "v10" && #defaultGraphPlugin) || (#graph != _|_ && #grafanaVersion == "v7" && !#defaultGraphPlugin) {
		type: _ | *"timeseries"

		options: legend: displayMode: _ | *#graph.legend.type
		options: legend: placement:   _ | *#graph.legend.placement
		options: legend: calcs:       _ | *#graph.legend.values

		fieldConfig: defaults: color: mode: "palette-classic" // for some reason CUE refuse to trim this value if it will be default...
		if #graph.display.bars.use {
			fieldConfig: defaults: custom: drawStyle: _ | *"bars"
			fieldConfig: defaults: custom: lineWidth: _ | *1
			fieldConfig: defaults: custom: pointSize: _ | *5
		}
		if !#graph.display.bars.use && #graph.display.points.use {
			fieldConfig: defaults: custom: drawStyle: _ | *"points"
			fieldConfig: defaults: custom: lineWidth: _ | *1
			fieldConfig: defaults: custom: pointSize: _ | *#graph.display.points.size
		}
		if !#graph.display.bars.use && !#graph.display.points.use {
			fieldConfig: defaults: custom: drawStyle: _ | *"lines"
			fieldConfig: defaults: custom: lineWidth: _ | *#graph.display.lines.size
			fieldConfig: defaults: custom: pointSize: _ | *5
		}
		if #graph.display.nulls == "connected" {
			fieldConfig: defaults: custom: insertNulls: _ | *false
		}

		fieldConfig: defaults: custom: fillOpacity:   _ | *0
		fieldConfig: defaults: custom: axisLabel:     _ | *""
		fieldConfig: defaults: custom: axisPlacement: _ | *"auto"
		fieldConfig: defaults: custom: barAlignment:  _ | *0
		fieldConfig: defaults: custom: gradientMode:  _ | *"none"
		fieldConfig: defaults: custom: hideFrom: graph:   _ | *false
		fieldConfig: defaults: custom: hideFrom: legend:  _ | *false
		fieldConfig: defaults: custom: hideFrom: tooltip: _ | *false
		fieldConfig: defaults: custom: lineInterpolation: _ | *"linear"
		fieldConfig: defaults: custom: scaleDistribution: type: _ | *"linear"
		fieldConfig: defaults: custom: showPoints: _ | *"never"
		fieldConfig: defaults: custom: spanNulls:  _ | *true

		if #grafanaVersion == "v10" {
			options: legend: showLegend: _ | *true
			options: tooltip: mode:      _ | *"multi"
			options: tooltip: sort:      _ | *"desc"
			fieldConfig: defaults: custom: axisCenteredZero: _ | *false
			fieldConfig: defaults: custom: axisColorMode:    _ | *"text"
			fieldConfig: defaults: custom: hideFrom: legend:      _ | *false
			fieldConfig: defaults: custom: hideFrom: tooltip:     _ | *false
			fieldConfig: defaults: custom: hideFrom: viz:         _ | *false
			fieldConfig: defaults: custom: stacking: group:       _ | *"A"
			fieldConfig: defaults: custom: stacking: mode:        _ | *"none"
			fieldConfig: defaults: custom: thresholdsStyle: mode: _ | *"off"
			fieldConfig: defaults: unit: _ | *#graph.leftY.unit
			fieldConfig: defaults: unitScale: _ | *true

		}

		options: tooltipOptions: mode: _ | *"single"

		yaxes: [for i, axis in yaxes {
			if i == 0 {
				format: _ | *#graph.leftY.unit
				if #graph.leftY.hashKey != _|_ {$$hashKey: _ | *#graph.leftY.hashKey}
			}
			if i == 1 {
				format: _ | *#graph.rightY.unit
				if #graph.rightY.hashKey != _|_ {$$hashKey: _ | *#graph.rightY.hashKey}
			}
			logBase: _ | *1
			show:    _ | *true
		}]
	}

	if (#graph != _|_ && #grafanaVersion == "v7" && #defaultGraphPlugin) || (#graph != _|_ && #grafanaVersion == "v10" && !#defaultGraphPlugin) {
		type:        _ | *"graph"
		description: _ | *""
		xaxis: mode:         _ | *"time"
		xaxis: show:         _ | *true
		tooltip: shared:     _ | *true
		tooltip: value_type: _ | *"individual"
		lines:       _ | *#graph.display.lines.use
		points:      _ | *#graph.display.points.use
		bars:        _ | *#graph.display.bars.use
		linewidth:   _ | *#graph.display.lines.size
		pointradius: _ | *#graph.display.points.size

		dashes:        _ | *false
		dashLength:    _ | *10
		fillGradient:  _ | *0
		hiddenSeries:  _ | *false
		nullPointMode: _ | *#graph.display.nulls
		percentage:    _ | *false
		pluginVersion: _ | *#pluginVersion[type]
		renderer:      _ | *"flot"
		spaceLength:   _ | *10
		stack:         _ | *#graph.display.stack

		fill: _ | *#graph.display.fill
		tooltip: sort:   _ | *#graph.tooltip.sort
		legend: avg:     _ | *list.Contains(#graph.legend.values, "mean")
		legend: max:     _ | *list.Contains(#graph.legend.values, "max")
		legend: min:     _ | *list.Contains(#graph.legend.values, "min")
		legend: current: _ | *list.Contains(#graph.legend.values, "lastNotNull")
		legend: total:   _ | *list.Contains(#graph.legend.values, "sum")

		legend: rightSide: _ | *(#graph.legend.placement == "right")
		if (#graph.legend.sortBy & string) != _|_ {
			#v7ToV10: {min: "min", max: "max", mean: "avg", avg: "avg", total: "total", sum: "total", current: "current", lastNotNull: "current", last: "current"}
			legend: sort: _ | *#v7ToV10[#graph.legend.sortBy]
		}
		legend: sortDesc: _ | *(#graph.legend.sortHow == "desc")

		legend: alignAsTable: _ | *(#graph.legend.type == "table")
		legend: values:       _ | *#graph.legend.showValues
		legend: show:         _ | *true

		yaxes: [for i, axis in yaxes {
			if i == 0 {
				format: _ | *#graph.leftY.unit
				if #graph.leftY.hashKey != _|_ {$$hashKey: _ | *#graph.leftY.hashKey}
			}
			if i == 1 {
				format: _ | *#graph.rightY.unit
				if #graph.rightY.hashKey != _|_ {$$hashKey: _ | *#graph.rightY.hashKey}
			}
			logBase: _ | *1
			show:    _ | *true
		}]
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
			expr: _ | *#prom
			if #grafanaVersion == "v7" {interval: _ | *""}
		}
		#mqlVisual?: {
			aliasBy:     string
			projectName: string
			metricType:  string
			filters: [...string]
			groupBys: [...string]
			unit:               string
			alignmentPeriod:    string | *"cloud-monitoring-auto"
			crossSeriesReducer: *"REDUCE_MEAN" | "REDUCE_SUM" | "REDUCE_MAX"
			perSeriesAligner:   *"ALIGN_MEAN" | "ALIGN_INTERPOLATE" | "ALIGN_NEXT_OLDER" | "ALIGN_RATE"
			metricKind:         "CUMULATIVE" | *"GAUGE"
			valueType:          "INT64" | "DOUBLE"
			sloV7:              bool | *false
		}
		if #mqlVisual != _|_ {
			metricQuery: aliasBy:            _ | *#mqlVisual.aliasBy
			metricQuery: alignmentPeriod:    _ | *#mqlVisual.alignmentPeriod
			metricQuery: crossSeriesReducer: _ | *#mqlVisual.crossSeriesReducer
			metricQuery: filters:            _ | *#mqlVisual.filters
			metricQuery: groupBys:           _ | *#mqlVisual.groupBys
			metricQuery: metricKind:         _ | *#mqlVisual.metricKind
			metricQuery: metricType:         _ | *#mqlVisual.metricType
			metricQuery: perSeriesAligner:   _ | *#mqlVisual.perSeriesAligner
			metricQuery: projectName:        _ | *#mqlVisual.projectName
			metricQuery: unit:               _ | *#mqlVisual.unit
			metricQuery: valueType:          _ | *#mqlVisual.valueType
			metricQuery: query:              _ | *""
			metricQuery: editorMode:         _ | *"visual"
			if #grafanaVersion == "v10" {metricQuery: preprocessor: _ | *"none"}
			if #grafanaVersion == "v7" && #mqlVisual.sloV7 {
				sloQuery: {
					projectName:     _ | *#mqlVisual.projectName
					alignmentPeriod: _ | *"cloud-monitoring-auto"
					selectorName:    _ | *"select_slo_health"
					aliasBy:         _ | *""
					serviceId:       _ | *""
					serviceName:     _ | *""
					sloId:           _ | *""
					sloName:         _ | *""
				}
			}
		}

		#mqlScript?: {
			query:              string
			aliasBy:            string
			projectName:        string
			alignmentPeriod:    string | *"cloud-monitoring-auto"
			crossSeriesReducer: *"REDUCE_MEAN" | "REDUCE_SUM" | "REDUCE_MAX"
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
		if #grafanaVersion == "v7" && (type == "stat" || type == "table") {
			exemplar: _ | *false
			instant:  _ | *true
		}
		if #grafanaVersion == "v10" && (type == "stat" || type == "table") {
			exemplar: _ | *false
			instant:  _ | *false
		}
		if type == "timeseries" || type == "graph" {
			instant:  _ | *false
			exemplar: _ | *true
			//			range:    _ | *true
		}
	}]
}

#monPanel: {
	#alert?: _

	if #grafanaVersion == "v7" && #alert != _|_ {
		options: alertThreshold: _
		alert: {for key, value in #alert {"\(key)": value}}

		if options.alertThreshold {
			thresholds: [for i, threshold in thresholds {
				colorMode: _ | *"critical"
				fill:      _ | *true
				line:      _ | *true
				if #alert.conditions[0].evaluator.type == "outside_range" && i == 0 {op: _ | *"lt"}
				if #alert.conditions[0].evaluator.type == "outside_range" && i == 1 {op: _ | *"gt"}
				if #alert.conditions[0].evaluator.type != "outside_range" {op: _ | *#alert.conditions[0].evaluator.type}
				value:   _ | *#alert.conditions[0].evaluator.params[i]
				visible: _ | *true
			}]
		}
	}
}

#monPanel: {
	panelH=h?: number
	if panelH != _|_ {gridPos: h: _ | *panelH}

	type:          string
	pluginVersion: _ | *#pluginVersion[type]

	if type == "stat" {
		options: reduceOptions: fields: _ | *""
		options: reduceOptions: values: _ | *false

		options: colorMode:   _ | *"value"
		options: graphMode:   _ | *"area"
		options: justifyMode: _ | *"auto"
		options: orientation: _ | *"auto"
		options: textMode:    _ | *"auto"

		if #grafanaVersion == "v7" {
			options: alertThreshold: _ | *true
			description:   _ | *""
			timeShift:     _ | *null
			timeFrom:      _ | *null
			steppedLine:   _ | *false
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
			yaxis: align:        _ | *false
			xaxis: mode:         _ | *"time"
			xaxis: show:         _ | *true
		}
	}
	if type == "graph" {
		steppedLine: _ | *false
		yaxis: align:            _ | *false
		options: alertThreshold: _ | *true
	}
}
