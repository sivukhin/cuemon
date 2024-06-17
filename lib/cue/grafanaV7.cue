package cuemon

#schemaVersion:  27
#grafanaVersion: "v7"
#pluginVersion: {
	text:       "7.5.17"
	stat:       "7.5.17"
	graph:      "7.5.17"
	timeseries: "7.5.17"
}

#grafana: {
	gnetId?:       null
	style:         string
	timezone:      string
	title:         string
	uid:           string
	description?:  string | null
	graphTooltip:  number
	id:            number
	schemaVersion: number
	iteration?:    number | null
	version?:      number | null
	editable:      bool
	refresh?:      string | bool | null
	tags: [...string]
	time: {
		from: string
		to:   string
	}
	timepicker: refresh_intervals?: [...string] | null
	annotations: list: [...#annotation]
	links: [...#link]
	templating: list: [...#template]
	panels: [...#panel]
}
#grid: {
	h: number
	w: number
	x: number
	y: number
} | null
#annotation: {
	datasource: string
	iconColor:  string
	name:       string
	type:       string
	builtIn?:   number | null
	limit?:     number | null
	showIn?:    number | null
	enable:     bool
	hide:       bool
	tags?: [...string] | null
}
#datasource: string | null
#target: {
	refId:           string
	datasource?:     #datasource
	expr?:           string | null
	format?:         string | null
	interval?:       string | null
	legendFormat?:   string | null
	queryType?:      string | null
	intervalFactor?: number | null
	panelId?:        number | null
	exemplar?:       bool | null
	hide?:           bool | null
	instant?:        bool | null
	sloQuery?: {
		aliasBy?:         string | null
		alignmentPeriod?: string | null
		projectName?:     string | null
		selectorName?:    string | null
		serviceId?:       string | null
		serviceName?:     string | null
		sloId?:           string | null
		sloName?:         string | null
	} | null
	metricQuery?: {
		aliasBy:            string
		alignmentPeriod:    string
		crossSeriesReducer: string
		editorMode:         string
		metricKind:         string
		metricType:         string
		perSeriesAligner:   string
		projectName:        string
		query:              string
		unit:               string
		valueType:          string
		filters: [...string]
		groupBys: [...string]
	} | null
}
#template: {
	#_constant: {
		description?: null
		error?:       null
		name:         string
		query:        string
		datasource?:  #datasource
		label?:       string | null
		hide:         number
		skipUrlSync:  bool
	}
	#_custom: {
		allValue?:    null
		error?:       null
		name:         string
		query:        string
		datasource?:  #datasource
		description?: string | null
		label?:       string | null
		queryValue?:  string | null
		hide:         number
		includeAll:   bool
		multi:        bool
		skipUrlSync:  bool
		tags?: [] | null
		options: [...{
			text:     string
			value:    string
			selected: bool
		}]
		current: {
			selected: bool
			text: [...string] | string
			value: [...string] | string
		}
	}
	#_query: {
		description?:   null
		error?:         null
		definition:     string
		name:           string
		regex:          string
		tagValuesQuery: string
		tagsQuery:      string
		allValue?:      string | null
		datasource?:    #datasource
		label?:         string | null
		hide:           number
		refresh:        number
		sort:           number
		includeAll:     bool
		multi:          bool
		skipUrlSync:    bool
		useTags:        bool
		tags: []
		options: [...{
			text:     string
			value:    string
			selected: bool
		}]
		current: {
			selected: bool
			tags?: [] | null
			text: [...string] | string
			value: [...string] | string
		}
		query: {
			refId:               string
			labelKey?:           string | null
			projectName?:        string | null
			query?:              string | null
			selectedMetricType?: string | null
			selectedQueryType?:  string | null
			selectedSLOService?: string | null
			selectedService?:    string | null
			loading?:            bool | null
			sloServices?: [] | null
			projects?: [...{
				label: string
				value: string
			}] | null
		}
	}
	#_textbox: {
		description?: null
		error?:       null
		label?:       null
		name:         string
		query:        string
		datasource?:  #datasource
		hide:         number
		skipUrlSync:  bool
		current: {
			text:     string
			value:    string
			selected: bool
		}
		options: [...{
			text:     string
			value:    string
			selected: bool
		}]
	}
	type:        "constant" | "custom" | "query" | "textbox"
	datasource?: #datasource
	if type == "constant" {
		#_constant
	}
	if type == "custom" {
		#_custom
	}
	if type == "query" {
		#_query
	}
	if type == "textbox" {
		#_textbox
	}
}
#yaxes: {
	format:     string
	$$hashKey?: string | null
	label?:     string | null
	max?:       string | null
	min?:       string | null
	logBase:    number
	show:       bool
}
#seriesOverrides: {
	alias:         string
	$$hashKey?:    string | null
	color?:        string | null
	fill?:         number | null
	fillGradient?: number | null
	linewidth?:    number | null
	yaxis?:        number | null
	zindex?:       number | null
	dashes?:       bool | null
	hiddenSeries?: bool | null
	legend?:       bool | null
}
#link: {
	#_dashboards: {
		icon:        string
		title:       string
		tooltip:     string
		url:         string
		asDropdown:  bool
		includeVars: bool
		keepTime:    bool
		targetBlank: bool
		tags: [...string]
	}
	#_link: {
		icon:         string
		title:        string
		tooltip:      string
		url:          string
		targetBlank:  bool
		asDropdown?:  bool | null
		includeVars?: bool | null
		keepTime?:    bool | null
		tags: []
	}
	type: "dashboards" | "link"
	if type == "dashboards" {
		#_dashboards
	}
	if type == "link" {
		#_link
	}
}
#panel: {
	#_gauge: {
		timeFrom?:     null
		timeShift?:    null
		description:   string
		pluginVersion: string
		title:         string
		datasource?:   #datasource
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		yaxes?: [...#yaxes] | null
		gridPos?: #grid
		options: {
			orientation:          string
			showThresholdLabels:  bool
			showThresholdMarkers: bool
			text: titleSize: number
			reduceOptions: {
				fields: string
				values: bool
				calcs: [...string]
			}
		}
		fieldConfig: {
			overrides: []
			defaults: {
				unit: string
				mappings: []
				color: mode: string
				links: [...{
					title:       string
					url:         string
					targetBlank: bool
				}]
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
			}
		}
	}
	#_graph: {
		maxDataPoints?: null
		timeFrom?:      null
		timeShift?:     null
		nullPointMode:  string
		pluginVersion:  string
		renderer:       string
		title:          string
		datasource?:    #datasource
		description?:   string | null
		height?:        string | null
		interval?:      string | null
		dashLength:     number
		fill:           number
		fillGradient:   number
		linewidth:      number
		pointradius:    number
		spaceLength:    number
		decimals?:      number | null
		id?:            number | null
		bars:           bool
		dashes:         bool
		hiddenSeries:   bool
		lines:          bool
		percentage:     bool
		points:         bool
		stack:          bool
		steppedLine:    bool
		editable?:      bool | null
		error?:         bool | null
		isNew?:         bool | null
		targets: [...#target]
		timeRegions: []
		links?: [] | null
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		yaxes?: [...#yaxes] | null
		grid?: {} | null
		aliasColors: {
			"max - istio-proxy"?: string | null
			"max - master"?:      string | null
		}
		yaxis: {
			alignLevel?: null
			align:       bool
		}
		gridPos?: #grid
		tooltip: {
			value_type:    string
			sort:          number
			shared:        bool
			msResolution?: bool | null
		}
		options: {
			alertThreshold: bool
			sortBy?: [] | null
		}
		xaxis: {
			buckets?: null
			name?:    null
			mode:     string
			show:     bool
			values: []
		}
		transformations?: [...{
			id: string
			options: {}
		}] | null
		thresholds: [...{
			colorMode:  string
			op:         string
			$$hashKey?: string | null
			yaxis?:     string | null
			value?:     number | null
			fill:       bool
			line:       bool
			visible?:   bool | null
		}]
		legend: {
			sideWidth?:    null
			sort?:         string | null
			avg:           bool
			current:       bool
			max:           bool
			min:           bool
			show:          bool
			total:         bool
			values:        bool
			alignAsTable?: bool | null
			hideEmpty?:    bool | null
			hideZero?:     bool | null
			rightSide?:    bool | null
			sortDesc?:     bool | null
		}
		fieldConfig: {
			defaults: {
				unit?: string | null
				links?: [...{
					title:       string
					url:         string
					targetBlank: bool
				}] | null
			}
			overrides: [...{
				matcher: {
					id:      string
					options: string
				}
				properties: [...{
					id: string
					value: [...{
						title:       string
						url:         string
						targetBlank: bool
					}]
				}]
			}]
		}
		alert?: {
			executionErrorState: string
			for:                 string
			frequency:           string
			name:                string
			noDataState:         string
			message?:            string | null
			handler:             number
			alertRuleTags: {
				group?:       string | null
				og_priority?: string | null
				project?:     string | null
			}
			notifications: [...{
				uid: string
			}]
			conditions: [...{
				type: string
				operator: type: string
				query: params: [...string]
				reducer: {
					type: string
					params: []
				}
				evaluator: {
					type: string
					params: [...number]
				}
			}]
		} | null
	}
	#_row: {
		title:       string
		datasource?: #datasource
		id?:         number | null
		collapsed:   bool
		panels: [...#panel]
		seriesOverrides?: [...#seriesOverrides] | null
		targets?: [...#target] | null
		yaxes?: [...#yaxes] | null
		gridPos?: #grid
	}
	#_stat: {
		timeFrom?:      null
		timeShift?:     null
		pluginVersion:  string
		title:          string
		datasource?:    #datasource
		description?:   string | null
		interval?:      string | null
		nullPointMode?: string | null
		renderer?:      string | null
		dashLength?:    number | null
		fill?:          number | null
		fillGradient?:  number | null
		id?:            number | null
		linewidth?:     number | null
		pointradius?:   number | null
		spaceLength?:   number | null
		bars?:          bool | null
		dashes?:        bool | null
		hiddenSeries?:  bool | null
		lines?:         bool | null
		percentage?:    bool | null
		points?:        bool | null
		stack?:         bool | null
		steppedLine?:   bool | null
		targets: [...#target]
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		thresholds?: [] | null
		timeRegions?: [] | null
		yaxes?: [...#yaxes] | null
		aliasColors?: {} | null
		yaxis?: {
			alignLevel?: null
			align:       bool
		} | null
		tooltip?: {
			value_type: string
			sort:       number
			shared:     bool
		} | null
		gridPos?: #grid
		xaxis?: {
			buckets?: null
			name?:    null
			mode:     string
			show:     bool
			values: []
		} | null
		legend?: {
			sort?:        null
			alignAsTable: bool
			avg:          bool
			current:      bool
			max:          bool
			min:          bool
			rightSide:    bool
			show:         bool
			sortDesc:     bool
			total:        bool
			values:       bool
		} | null
		options: {
			graphMode:       string
			textMode:        string
			colorMode?:      string | null
			justifyMode?:    string | null
			orientation?:    string | null
			alertThreshold?: bool | null
			sortBy?: [] | null
			text?: {} | null
			reduceOptions: {
				fields: string
				values: bool
				calcs: [...string]
			}
		}
		fieldConfig: {
			overrides: []
			defaults: {
				unit?: string | null
				min?:  number | null
				mappings?: [] | null
				color?: {
					mode: string
				} | null
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
			}
		}
	}
	#_table: {
		timeFrom?:     null
		timeShift?:    null
		pluginVersion: string
		title:         string
		datasource?:   #datasource
		id?:           number | null
		targets: [...#target]
		links?: [] | null
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		yaxes?: [...#yaxes] | null
		gridPos?: #grid
		options: {
			frameIndex?: number | null
			showHeader:  bool
			sortBy?: [...{
				displayName: string
				desc:        bool
			}] | null
		}
		transformations: [...{
			id: string
			options: {
				reducers?: [...string] | null
				excludeByName?: {} | null
				renameByName?: {
					Field?:                     string | null
					Max?:                       string | null
					Total?:                     string | null
					"variant (distinctCount)"?: string | null
				} | null
				indexByName?: {
					cg_name?:                   number | null
					cg_type?:                   number | null
					"variant (distinctCount)"?: number | null
				} | null
				sort?: [...{
					field: string
					desc:  bool
				}] | null
				fields?: {
					cg_name?: {
						operation: string
						aggregations: []
					} | null
					cg_type?: {
						operation: string
						aggregations: []
					} | null
					variant?: {
						operation: string
						aggregations: [...string]
					} | null
				} | null
			}
		}]
		fieldConfig: {
			overrides: [...{
				matcher: {
					id:      string
					options: string
				}
				properties: [...{
					id: string
					value: [...{
						title:       string
						url:         string
						targetBlank: bool
					}] | string | bool
				}]
			}]
			defaults: {
				unit?: string | null
				mappings: []
				color: mode: string
				custom: {
					align?:       string | null
					displayMode?: string | null
					filterable:   bool
				}
				links?: [...{
					title:       string
					url:         string
					targetBlank: bool
				}] | null
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
			}
		}
	}
	#_text: {
		timeFrom?:     null
		timeShift?:    null
		pluginVersion: string
		title:         string
		datasource?:   #datasource
		id?:           number | null
		transparent:   bool
		targets: [...#target]
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		yaxes?: [...#yaxes] | null
		options: {
			content: string
			mode:    string
		}
		gridPos?: #grid
		fieldConfig: {
			overrides: []
			defaults: {}
		}
	}
	#_timeseries: {
		maxDataPoints?: null
		timeFrom?:      null
		timeShift?:     null
		description:    string
		interval:       string
		pluginVersion:  string
		title:          string
		datasource?:    #datasource
		id?:            number | null
		targets: [...#target]
		panels?: [...#panel] | null
		seriesOverrides?: [...#seriesOverrides] | null
		yaxes?: [...#yaxes] | null
		gridPos?: #grid
		options: {
			graph: {}
			tooltipOptions: mode: string
			legend: {
				displayMode: string
				placement:   string
				calcs: [...string]
			}
		}
		fieldConfig: {
			overrides: [...{
				"__systemRef": string
				matcher: {
					id: string
					options: {
						mode:     string
						prefix:   string
						readOnly: bool
						names: [...string]
					}
				}
				properties: [...{
					id: string
					value: {
						graph:   bool
						legend:  bool
						tooltip: bool
					}
				}]
			}]
			defaults: {
				unit: string
				min:  number
				mappings: []
				color: mode: string
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
				custom: {
					axisLabel:         string
					axisPlacement:     string
					drawStyle:         string
					gradientMode:      string
					lineInterpolation: string
					showPoints:        string
					barAlignment:      number
					fillOpacity:       number
					lineWidth:         number
					pointSize:         number
					spanNulls:         bool
					scaleDistribution: type: string
					hideFrom: {
						graph:   bool
						legend:  bool
						tooltip: bool
					}
				}
			}
		}
	}
	type:        "gauge" | "graph" | "row" | "stat" | "table" | "text" | "timeseries"
	datasource?: #datasource
	id?:         number
	panels: [...#panel]
	gridPos?: #grid
	yaxes?: [...#yaxes]
	seriesOverrides?: [...#seriesOverrides]
	targets: [...#target]
	if type == "gauge" {
		#_gauge
	}
	if type == "graph" {
		#_graph
	}
	if type == "row" {
		#_row
	}
	if type == "stat" {
		#_stat
	}
	if type == "table" {
		#_table
	}
	if type == "text" {
		#_text
	}
	if type == "timeseries" {
		#_timeseries
	}
}
