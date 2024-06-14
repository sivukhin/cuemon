package cuemon

import (
	"strings"
)

#monVariable: #template & {#template}
#monVariable: {
	#hideTypes: {none: 0, label: 1, variable: 2}
	#refreshTypes: {none: 0, dashboardLoad: 1, timeRangeChange: 2}
	#sortTypes: {none: 0, alphAsc: 1, alphDesc: 2, numAsc: 3, numDesc: 4, ialphAsc: 5, ialphDesc: 6}
	#queryTypes: {labelNames: 0, labelValues: 1, metrics: 2, queryResult: 3, seriesQuery: 4, classicQuery: 5}

	#datasrc?: _
	#constant?: {
		name:  string
		value: string
	}
	#custom?: {
		name:  string
		label: string | null | *name
		hide:  *"none" | "label" | "variable"
		multi: bool | *true
		all:   bool | *multi
		options: [...{value: string, text: string | *value, selected: bool | *false}]
	}
	#query?: {
		name:     string
		label:    string | null | *name
		hide:     *"none" | "label" | "variable"
		refresh:  *"dashboardLoad" | "timeRangeChange"
		sort:     "none" | "alphAsc" | "alphDesc" | "numAsc" | "numDesc" | *"ialphAsc" | "ialphDesc"
		multi:    bool | *true
		all:      bool | *multi
		allValue?: string
		regex:    string | *""
		query:    string
		current: [...string] | string
	}
	#textbox?: {
		name:    string
		label:   string | *name
		default: string
		hide:    *"none" | "label" | "variable"
	}

	if #datasrc != _|_ {
		datasource: #datasrc
	}
	if #constant != _|_ {
		type:        "constant"
		name:        _ | *#constant.name
		query:       _ | *#constant.value
		skipUrlSync: _ | *false
		hide:        _ | *2
		if #grafanaVersion == "v7" {
			label: _ | *#constant.name
		}
	}
	if #custom != _|_ {
		type:        "custom"
		label:       _ | *#custom.label
		name:        _ | *#custom.name
		multi:       _ | *#custom.multi
		includeAll:  _ | *#custom.all
		options:     _ | *#custom.options
		skipUrlSync: _ | *false
		hide:        _ | *#hideTypes[#custom.hide]
		query: _ | *strings.Join([for option in #custom.options {strings.Replace(option.value, ",", "\\,", -1)}], ",")
		current: selected: _ | *true
		if multi {
			current: text: _ | *[for option in #custom.options if option.selected {
				if option.text != "$__all" {option.text}
				if option.text == "$__all" {"All"}
			}]
			current: value: _ | *[for option in #custom.options if option.selected {option.value}]
		}
		if !multi {
			current: text:  _ | *[for option in #custom.options if option.selected {option.text}][0] | ""
			current: value: _ | *[for option in #custom.options if option.selected {option.value}][0] | ""
		}
		if #grafanaVersion == "v7" {description: string | null | *""}
//		queryValue: string | *""
	}
	if #query != _|_ {
		type:           "query"
		name:           _ | *#query.name
		label:          _ | *#query.label
		hide:           _ | *#hideTypes[#query.hide]
		includeAll:     _ | *#query.all
		multi:          _ | *#query.multi
		skipUrlSync:    _ | *false
		sort:           _ | *#sortTypes[#query.sort]
		regex:          _ | *#query.regex
		refresh:        _ | *#refreshTypes[#query.refresh]
		useTags:        _ | *false
		tagsQuery:      _ | *""
		tagValuesQuery: _ | *""

		definition: _ | *query.query
		query: query:      _ | *#query.query
		current: selected: _ | *true
		current: value:    _ | *#query.current
		if #query.multi {
			current: text: _ | *[for option in #query.current {
				if option != "$__all" {option}
				if option == "$__all" {"All"}
			}]
		}
		if !#query.multi {
			current: text: _ | *#query.current
		}
		if #query.allValue != _|_ {
			allValue: _ | *#query.allValue
		}
		if #grafanaVersion == "v7" {
			tags: _ | *[]
			query: refId: _ | *"StandardVariableQuery"
		}
		if #grafanaVersion == "v10" {
			query: refId: _ | *"PrometheusVariableQueryEditor-VariableQuery"
			if strings.HasPrefix(query.query, "label_values") {query: qryType: _ | *#queryTypes.labelValues}
			if strings.HasPrefix(query.query, "query_result") {query: qryType: _ | *#queryTypes.queryResult}
		}
	}
	if #textbox != _|_ {
		type:        "textbox"
		name:        _ | *#textbox.name
		label:       _ | *#textbox.label
		query:       _ | *#textbox.default
		hide:        _ | *#hideTypes[#textbox.hide]
		skipUrlSync: _ | *false
		options: _ | *[current]
		current: selected: _ | *false
		current: text:     _ | *#textbox.default
		current: value:    _ | *#textbox.default
	}
}
