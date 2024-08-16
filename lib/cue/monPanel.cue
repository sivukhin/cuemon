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
		html?:       string
		markdown?:   string
		transparent: bool | *false
	}
	if #text != _|_ {
		type:          "text"
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
		options: code: language:        _ | *"plaintext"
		options: code: showLineNumbers: _ | *false
		options: code: showMiniMap:     _ | *false
		error:    _ | *false
		editable: _ | *true
	}
}

#monPanel: {
	#stat?: {
		reducer: [...("lastNotNull" | "last" | "min" | "max" | "avg" | "mean")]
		yPrimary: unit: string | *"short"
		yPrimary: min?: number
		yPrimary: max?: number
		orientation: *"auto" | "vertical" | "horizontal"
	}

	if #stat != _|_ {
		type: "stat"
		options: orientation: _ | *#stat.orientation
		options: reduceOptions: calcs: _ | *#stat.reducer
		fieldConfig: defaults: unit:   _ | *#stat.yPrimary.unit
		if #stat.yPrimary.min != _|_ {fieldConfig: defaults: min: _ | *#stat.yPrimary.min}
		if #stat.yPrimary.max != _|_ {fieldConfig: defaults: max: _ | *#stat.yPrimary.max}
	}
}

#monPanel: {
	type: string
	targets: [...]
	#d=datasource: #datasource

	#target: {
		prom?: string
		mqlVisual?: {
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
		mqlScript?: {
			query:              string
			aliasBy:            string
			projectName:        string
			alignmentPeriod:    string | *"cloud-monitoring-auto"
			crossSeriesReducer: *"REDUCE_MEAN" | "REDUCE_SUM" | "REDUCE_MAX"
			perSeriesAligner:   *"ALIGN_MEAN" | "ALIGN_INTERPOLATE" | "ALIGN_NEXT_OLDER" | "ALIGN_RATE"
		}
		format: string | *""
		hide:   bool | *false
	}
	#targets: [...#target]
	targets: [for _ in #targets { _ }]
	targets: [for i, target in targets {
		target

		datasource?: _ | *#d
		if #targets[i].prom != _|_ {
			expr: _ | *#targets[i].prom
		}
		if #targets[i].mqlVisual != _|_ {
			metricQuery: aliasBy:            _ | *#targets[i].mqlVisual.aliasBy
			metricQuery: alignmentPeriod:    _ | *#targets[i].mqlVisual.alignmentPeriod
			metricQuery: crossSeriesReducer: _ | *#targets[i].mqlVisual.crossSeriesReducer
			metricQuery: filters:            _ | *#targets[i].mqlVisual.filters
			metricQuery: groupBys:           _ | *#targets[i].mqlVisual.groupBys
			metricQuery: metricKind:         _ | *#targets[i].mqlVisual.metricKind
			metricQuery: metricType:         _ | *#targets[i].mqlVisual.metricType
			metricQuery: perSeriesAligner:   _ | *#targets[i].mqlVisual.perSeriesAligner
			metricQuery: projectName:        _ | *#targets[i].mqlVisual.projectName
			metricQuery: unit:               _ | *#targets[i].mqlVisual.unit
			metricQuery: valueType:          _ | *#targets[i].mqlVisual.valueType
			metricQuery: query:              _ | *""
			metricQuery: editorMode:         _ | *"visual"
			metricQuery: preprocessor:       _ | *"none"
		}
		if #targets[i].mqlScript != _|_ {
			metricQuery: aliasBy:            _ | *#targets[i].mqlScript.aliasBy
			metricQuery: alignmentPeriod:    _ | *#targets[i].mqlScript.alignmentPeriod
			metricQuery: crossSeriesReducer: _ | *#targets[i].mqlScript.crossSeriesReducer
			metricQuery: perSeriesAligner:   _ | *#targets[i].mqlScript.perSeriesAligner
			metricQuery: projectName:        _ | *#targets[i].mqlScript.projectName
			metricQuery: query:              _ | *#targets[i].mqlScript.query

			metricQuery: editorMode: _ | *"mql"
			metricQuery: metricKind: _ | *""
			metricQuery: metricType: _ | *""
			metricQuery: unit:       _ | *""
			metricQuery: valueType:  _ | *""
		}
		queryType:    _ | *"metrics"
		legendFormat: _ | *#targets[i].format
		hide:         _ | *#targets[i].hide

		if type == "stat" || type == "table" {
			exemplar: _ | *false
			instant:  _ | *false
		}
		if type == "timeseries" {
			instant:  _ | *false
			exemplar: _ | *true
			//			range:    _ | *true
		}
	}]
}

#monPanel: {
	#alphabet: [for x in list.Range(10, 256+10, 1) {strings.ToUpper(strconv.FormatInt(x, 36))}]
	#targets: _

	targets: [for i, target in targets {
		refId: _ | *#alphabet[i]
	}]
}

#monPanel: {
	panelH=h?: number
	if panelH != _|_ {gridPos: h: _ | *panelH}

	type:          string
	pluginVersion: _ | *#pluginVersion[type]

	#title:       string
	#description: string | *""
	#datasrc?:    #datasource

	if type == "stat" || type == "timeseries" {
		title:       _ | *#title
		description: _ | *#description
		datasource:  _ | *#datasrc
	}

	if type == "stat" {
		options: reduceOptions: fields: _ | *""
		options: reduceOptions: values: _ | *false

		options: colorMode:   _ | *"value"
		options: graphMode:   _ | *"area"
		options: justifyMode: _ | *"auto"
		options: textMode:    _ | *"auto"
	}
}
