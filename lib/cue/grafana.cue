package cuemon

#GrafanaSchema: {
	#Tooltip: {Default: 0, SharedCrosshair: 1, SharedTooltip: 2}

	schemaVersion: number
	id:            int64
	uid:           string
	title:         string
	description?:  string
	iteration?:    int64
	gnetId:        string | *null
	style:         "light" | *"dark"
	timezone:      string | *"browser"
	editable:      bool | *true
	graphTooltip:  #Tooltip.Default | *#Tooltip.SharedCrosshair | #Tooltip.SharedTooltip
	time: {from: string | *"now-6h", to: string | *"now"}
	version?: int
	refresh:  string | *"15m"

	annotations: list:             [...#GrafanaAnnotation] | *[]
	timepicker: refresh_intervals: [...string] | *["10s", "30s", "1m", "5m", "15m", "30m", "1h"]

	tags?: [...string]
	links?: [...#GrafanaLink]
	panels?: [...{#GrafanaPanel, #schemaVersion: schemaVersion}]
	templating: list?: [...#GrafanaTemplate]
}

#GrafanaLink: {
	$$hashKey?:  string
	type:        "link" | "dashboards"
	icon:        *"external link" | "dashboard" | "question" | "info" | "bolt" | "doc" | "cloud"
	targetBlank: bool | *true
	if type == "link" {
		title:   string
		url:     string
		tooltip: string | *""
		tags: []
	}
	if type == "dashboards" {
		tags: [...string]
		asDropdown:  bool | *false
		includeVars: bool | *false
		keepTime:    bool | *false
	}
}

#Hide: {None: 0, Label: 1, Variable: 2}
#Sort: {Disabled: 0, AlphAsc: 1, AlphDesc: 2, NumAsc: 3, NumDesc: 4, AlphAscNonSensitive: 5, AlphDescNonSensitive: 6}

#GrafanaTemplate: {
	type:        "constant" | "custom" | "query"
	name:        string
	label:       string | *null
	description: string | *null
	hide:        #Hide.None | #Hide.Label | #Hide.Variable
	skipUrlSync: bool | *false
	error:       null
	if type == "constant" {
		query: string
		hide:  hide | *#Hide.Variable
	}
	if type == "custom" {
		query: string
		current: {
			selected: bool | *true
			tags: []
			if multi {
				text: [...or([ for option in options {option.text}])]
				value: [...or([ for option in options {option.value}])]
			}
			if !multi {
				text:  or([ for option in options {option.text}])
				value: or([ for option in options {option.value}])
			}
		}
		options: [...{selected: bool, text: string, value: string}]
		multi:      bool
		includeAll: bool | *multi
		allValue:   string | *null
		hide:       hide | *#Hide.None
		queryValue: ""
	}
	if type == "query" {
		datasource: string
		definition: string
		refresh:    number | *1
		multi:      bool
		includeAll: bool | *multi
		regex?:     string
		allValue:   string | *null
		sort:       *#Sort.Disabled | #Sort.AlphAsc | #Sort.AlphDesc | #Sort.NumAsc | #Sort.NumDesc | #Sort.AlphAscNonSensitive | #Sort.AlphDescNonSensitive
		current: {selected: bool | *true, tags: [], text: [...string], value: [...string]}
		query: {query: string | *definition, refId: string | *"StandardVariableQuery"}
		hide: hide | *#Hide.None
		options: []
		tags: []
		tagsQuery:      string | *""
		tagValuesQuery: string | *""
		useTags:        false
	}
}

#GrafanaAnnotation: {
	builtIn:    number
	datasource: string
	enable:     bool
	hide:       bool
	iconColor:  string
	name:       string
	type:       string
}

#GrafanaXAxis: {
	buckets: null
	mode:    string | *"time"
	name:    null
	show:    bool | *true
	values: []
}

#GrafanaYAxis: {
	align:      bool | *false
	alignLevel: null
}

#GrafanaYAxes: {
	$$hashKey?: string
	format:     string
	logBase:    number | *1
	max:        string | *null
	min:        string | *null
	show:       bool | *true
	label:      string | *null
}

#GrafanaThreshold: {
	$$hashKey?: string
	colorMode:  string | *"critical"
	fill:       bool | *true
	line:       bool | *true
	op:         string | *"gt"
	value:      number
	visible:    bool | *true
	yaxis?:     string
}

#GrafanaLegend: {
	#schemaVersion: number
	show:           bool | *false
	values:         bool | *false
	current:        bool | *false
	avg:            bool | *false
	max:            bool | *false
	min:            bool | *false
	total:          bool | *false
	sort:           "current" | "avg" | "max" | "min" | "total" | *null
	sortDesc:       bool | *false
	rightSide:      bool | *false
	alignAsTable:   bool | *false
	hideEmpty?:     bool
	hideZero?:      bool
}

#GrafanaAlert: {
	name:                string
	message:             string
	for:                 string
	frequency:           string
	executionErrorState: "alerting" | *"keep_state"
	noDataState:         "alerting" | "keep_state" | *"no_data" | "ok"
	handler:             number | *1
	alertRuleTags: [string]: string
	notifications: [...{uid: string}]
	conditions: [...{
		type: string | *"query"
		operator: type: string | *"and"
		query: params: [...string]
		reducer: {
			params: []
			type: string | *"avg"
		}
		evaluator: {
			type: string | "lt" | "gt"
			params: [...number]
		}
	}]
}

#GrafanaTarget: {
	#schemaVersion: number
	refId:          string
	queryType?:     string
	hide:           bool | *false
	exemplar:       bool | *true
	interval:       string | *""
	if queryType == _|_ {
		expr:         string
		legendFormat: string
		format?:      string
		instant:      bool | *false
	}
	if queryType != _|_ {
		if queryType == "randomWalk" {
			expr:         string
			legendFormat: string
			format?:      string
			instant:      bool | *false
		}
		if queryType == "metrics" {
			metricQuery: {
				aliasBy:            string
				alignmentPeriod:    string
				crossSeriesReducer: string
				editorMode:         string
				metricKind:         string
				metricType:         string
				perSeriesAligner:   string
				projectName:        string
				unit:               string
				valueType:          string
				query:              string | *""
				preprocessor?:      string
				filters:            [...string] | *[]
				groupBys:           [...string] | *[]
			}
			sloQuery: {...}
		}
	}
}

#GrafanaTooltip: {
	#Sort: {None: 0, Increasing: 1, Decreasing: 2}
	shared:     bool | *true
	sort:       #Sort.None | #Sort.Increasing | #Sort.Decreasing
	value_type: string | *"individual"
}

#GrafanaPanel: {
	v=#schemaVersion: number
	id:               number
	title:            string
	type:             "row" | "graph" | "stat" | "table" | "timeseries"
	gridPos: {h: number, w: number, x: number, y: number}
	if type == "row" {
		collapsed: bool
		panels: [...{#GrafanaPanel, #schemaVersion: v}]
		datasource: null
	}
	if type == "graph" || type == "stat" || type == "table" {
		{#GrafanaGraph, #type: type, #schemaVersion: v}
	}
	if type == "timeseries" {
		{#GrafanaTimeseries, #schemaVersion: v}
	}
}

#GrafanaTimeseries: {
	v=#schemaVersion: number
	if #schemaVersion >= 37 {datasource: {type?: string, uid: string}}
	if #schemaVersion < 37 {datasource: string}
	options: {
		tooltip: {
			mode: "none" | "single" | *"multi"
			sort: "none" | "asc" | *"desc"
		}
		legend: {
			displayMode: "list" | *"table"
			placement:   "bottom" | *"right"
			showLegend:  bool | *true
			calcs: [...string]
		}
	}
	targets: [...{#GrafanaTarget, #schemaVersion: v}]
	fieldConfig: {
		overrides: [...]
		defaults: {
			mappings: []
			unit?: string
			color: mode: string | *"palette-classic"
			thresholds: {
				mode: string | *"absolute"
				steps: [...{color: string, value: number | null}]
			}
			custom: {
				axisCenteredZero:  bool | *false
				axisColorMode:     string | *"text"
				axisLabel:         string | *""
				axisPlacement:     string | *"auto"
				barAlignment:      number | *0
				drawStyle:         string | *"line"
				fillOpacity:       number | *0
				gradientMode:      string | *"none"
				lineInterpolation: string | *"linear"
				lineWidth:         number | *1
				pointSize:         number | *5
				showPoints:        string | *"auto"
				spanNulls:         bool | *false
				scaleDistribution: type: string | *"linear"
				thresholdsStyle: mode:   string | *"off"
				stacking: {
					group: string | *"A"
					mode:  string | *"none"
				}
				hideFrom: {
					legend:  bool | *false
					tooltip: bool | *false
					viz:     bool | *false
				}
			}
		}
	}
}

#GrafanaGraph: {
	v=#schemaVersion: number
	#type:            string
	datasource:       string
	description:      string | *""
	targets: [...{#GrafanaTarget, #schemaVersion: v}]
	tooltip?: #GrafanaTooltip
	thresholds: [...#GrafanaThreshold]
	yaxes: [...#GrafanaYAxes]
	legend: {#GrafanaLegend, #schemaVersion: v}
	alert?: #GrafanaAlert
	xaxis:  #GrafanaXAxis
	yaxis:  #GrafanaYAxis

	hiddenSeries:  bool | *false
	lines:         bool | *true
	nullPointMode: string | *"null"
	linewidth:     number | *1
	aliasColors?: {}
	bars:          bool | *false
	dashes:        bool | *false
	hiddenSeries:  bool | *false
	percentage:    bool | *false
	points:        bool | *false
	stack:         bool | *false
	steppedLine:   bool | *false
	dashLength:    number | *10
	spaceLength:   number | *10
	fill:          number | *0
	fillGradient:  number | *0
	pointradius:   number | *0
	description:   string | *""
	pluginVersion: string | *"7.5.17"
	renderer:      string | *"flot"
	timeFrom:      null
	timeRegions: []
	timeShift: null
	options: {
		alertThreshold: bool | *true
		colorMode?:     string
		graphMode?:     string
		justifyMode?:   string
		orientation?:   string
		textMode?:      string
		showHeader?:    bool
		reduceOptions?: {
			calcs: [...string]
			fields: string | *""
			values: bool
		}
		text?: {}
		sortBy: []
	}
	fieldConfig?: {
		overrides?: [
			...{
				matcher: {
					id:      string
					options: string
				}
				properties: [...{
					id:    string
					value: number
				}]
			},
		]
		defaults: {
			color?: mode: string
			custom?: {
				align:       string
				displayMode: string
				filterable:  bool
				width:       number
			}
			mappings?: []
			links?: []
			thresholds?: {
				mode: string | *"absolute"
				steps: [...{color: string, value: number | string | null}]
			}
		}
	}
	seriesOverrides?: [...{
		$$hashKey?:    string
		alias:         string
		color?:        string
		fill?:         number
		yaxis?:        1 | 2
		zindex?:       number
		dashes?:       bool
		hiddenSeries?: bool
		linewidth?:    number
	}]
	links?: [...{
		targetBlank: bool | *true
		title:       string
		url:         string
	}]
	transformations?: [...{...}]
}
