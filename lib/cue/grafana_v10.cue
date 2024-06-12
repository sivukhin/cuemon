package cuemon

#schemaVersion:  39
#grafanaVersion: "v10"
#pluginVersion: {
	text:  "10.3.0"
	table: "10.3.0"
}

#grafana: {
	refresh:              string
	timezone:             string
	title:                string
	uid:                  string
	weekStart:            string
	description?:         string | null
	fiscalYearStartMonth: number
	graphTooltip:         number
	id:                   number
	schemaVersion:        number
	iteration?:           number | null
	version?:             number | null
	editable:             bool
	liveNow:              bool
	tags: [...string]
	time: {
		from: string
		to:   string
	}
	timepicker: {
		refresh_intervals?: [...string] | null
		time_options?: [...string] | null
	}
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
	iconColor:    string
	name:         string
	expr?:        string | null
	step?:        string | null
	titleFormat?: string | null
	type?:        string | null
	builtIn?:     number | null
	enable:       bool
	hide:         bool
	datasource: {
		type: string
		uid:  string
	}
	target?: {
		type:     string
		limit:    number
		matchAny: bool
		tags: []
	} | null
}
#datasource: {
	type:  string
	uid:   string
	name?: string | null
} | string | null
#target: {
	refId:           string
	aliasBy?:        string | null
	editorMode?:     string | null
	expr?:           string | null
	expression?:     string | null
	format?:         string | null
	interval?:       string | null
	legendFormat?:   string | null
	metric?:         string | null
	queryType?:      string | null
	type?:           string | null
	intervalFactor?: number | null
	step?:           number | null
	exemplar?:       bool | null
	hide?:           bool | null
	instant?:        bool | null
	range?:          bool | null
	timeSeriesQuery?: {
		projectName: string
		query:       string
	} | null
	datasource?: #datasource
	sloQuery?: {
		aliasBy:          string
		alignmentPeriod:  string
		lookbackPeriod:   string
		perSeriesAligner: string
		projectName:      string
		selectorName:     string
		serviceId:        string
		serviceName:      string
		sloId:            string
		sloName:          string
	} | null
	metricQuery?: {
		projectName:         string
		aliasBy?:            string | null
		alignmentPeriod?:    string | null
		crossSeriesReducer?: string | null
		editorMode?:         string | null
		expr?:               string | null
		legendFormat?:       string | null
		metricKind?:         string | null
		metricType?:         string | null
		perSeriesAligner?:   string | null
		preprocessor?:       string | null
		query?:              string | null
		unit?:               string | null
		valueType?:          string | null
		range?:              bool | null
		filters?: [...string] | null
		groupBys?: [...string] | null
	} | null
}
#template: {
	#_constant: {
		name:        string
		query:       string
		hide:        number
		skipUrlSync: bool
		datasource?: #datasource
	}
	#_custom: {
		name:         string
		query:        string
		allFormat?:   string | null
		allValue?:    string | null
		label?:       string | null
		multiFormat?: string | null
		queryValue?:  string | null
		hide:         number
		refresh?:     number | null
		sort?:        number | null
		includeAll:   bool
		multi:        bool
		skipUrlSync:  bool
		datasource?:  #datasource
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
	#_datasource: {
		label:       string
		name:        string
		query:       string
		regex:       string
		queryValue?: string | null
		hide:        number
		refresh:     number
		includeAll:  bool
		multi:       bool
		skipUrlSync: bool
		options: []
		datasource?: #datasource
		current: {
			selected: bool
			text: [...string] | string
			value: [...string] | string
		}
	}
	#_interval: {
		auto_min:    string
		label:       string
		name:        string
		query:       string
		queryValue?: string | null
		auto_count:  number
		hide:        number
		refresh:     number
		auto:        bool
		skipUrlSync: bool
		datasource?: #datasource
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
	#_query: {
		definition:      string
		name:            string
		regex:           string
		allFormat?:      string | null
		allValue?:       string | null
		label?:          string | null
		multiFormat?:    string | null
		tagValuesQuery?: string | null
		tagsQuery?:      string | null
		hide:            number
		refresh:         number
		sort:            number
		includeAll:      bool
		multi:           bool
		skipUrlSync:     bool
		useTags?:        bool | null
		options: []
		query: {
			query:    string
			refId:    string
			qryType?: number | null
		}
		datasource?: #datasource
		current: {
			selected: bool
			isNone?:  bool | null
			text: [...string] | string
			value: [...string] | string
		}
	}
	#_textbox: {
		label:       string
		name:        string
		query:       string
		hide:        number
		skipUrlSync: bool
		datasource?: #datasource
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
	type:        "constant" | "custom" | "datasource" | "interval" | "query" | "textbox"
	datasource?: #datasource
	if type == "constant" {
		#_constant
	}
	if type == "custom" {
		#_custom
	}
	if type == "datasource" {
		#_datasource
	}
	if type == "interval" {
		#_interval
	}
	if type == "query" {
		#_query
	}
	if type == "textbox" {
		#_textbox
	}
}
#link: {
	#_dashboards: {
		icon:         string
		$$hashKey?:   string | null
		targetBlank:  bool
		asDropdown?:  bool | null
		includeVars?: bool | null
		keepTime?:    bool | null
		tags: [...string]
	}
	#_link: {
		icon:         string
		title:        string
		url:          string
		$$hashKey?:   string | null
		tooltip?:     string | null
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
	#_alertlist: {
		description: string
		title:       string
		id?:         number | null
		panels?: [...#panel] | null
		targets?: [...#target] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			alertInstanceLabelFilter: string
			alertName:                string
			groupMode:                string
			viewMode:                 string
			maxItems:                 number
			sortOrder:                number
			dashboardAlerts:          bool
			groupBy: []
			folder: {
				title: string
				id:    number
			}
			stateFilter: {
				error:   bool
				firing:  bool
				noData:  bool
				normal:  bool
				pending: bool
			}
		}
	}
	#_bargauge: {
		pluginVersion: string
		title:         string
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			displayMode:   string
			namePlacement: string
			orientation:   string
			sizing:        string
			valueMode:     string
			maxVizHeight:  number
			minVizHeight:  number
			minVizWidth:   number
			showUnfilled:  bool
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
				unit:      string
				max?:      number | null
				min?:      number | null
				unitScale: bool
				mappings: []
				color: mode: string
				thresholds: {
					mode: string
					steps: [...{
						color: string
						value: number | null
					}]
				}
			}
		}
	}
	#_dashlist: {
		pluginVersion: string
		title:         string
		id?:           number | null
		panels?: [...#panel] | null
		targets?: [...#target] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			folderUID:          string
			query:              string
			maxItems:           number
			includeVars:        bool
			keepTime:           bool
			showHeadings:       bool
			showRecentlyViewed: bool
			showSearch:         bool
			showStarred:        bool
			tags: [...string]
		}
	}
	#_gauge: {
		pluginVersion: string
		title:         string
		description?:  string | null
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			orientation:          string
			sizing:               string
			minVizHeight:         number
			minVizWidth:          number
			showThresholdLabels:  bool
			showThresholdMarkers: bool
			reduceOptions: {
				fields: string
				values: bool
				calcs: [...string]
			}
		}
		fieldConfig: {
			overrides: []
			defaults: {
				unit?:     string | null
				unitScale: bool
				mappings: []
				links?: [] | null
				color: mode: string
				thresholds: {
					mode: string
					steps: [...{
						color: string
						value: number | null
					}]
				}
			}
		}
	}
	#_graph: {
		nullPointMode:  string
		renderer:       string
		title:          string
		description?:   string | null
		interval?:      string | null
		pluginVersion?: string | null
		dashLength:     number
		fill:           number
		fillGradient:   number
		linewidth:      number
		pointradius:    number
		spaceLength:    number
		id?:            number | null
		span?:          number | null
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
		transparent?:   bool | null
		targets: [...#target]
		timeRegions: []
		panels?: [...#panel] | null
		yaxis: align: bool
		aliasColors: {
			Egress?:         string | null
			Ingress?:        string | null
			"test-index-1"?: string | null
		}
		datasource?: #datasource
		gridPos?:    #grid
		tooltip: {
			value_type:    string
			sort:          number
			shared:        bool
			msResolution?: bool | null
		}
		options?: {
			alertThreshold?: bool | null
			dataLinks?: [] | null
		} | null
		links?: [...{
			title:       string
			url:         string
			targetBlank: bool
		}] | null
		xaxis: {
			mode:     string
			format?:  string | null
			logBase?: number | null
			show:     bool
			values: []
		}
		yaxes: [...{
			format:     string
			$$hashKey?: string | null
			label?:     string | null
			max?:       string | null
			logBase:    number
			min?:       number | string | null
			show:       bool
		}]
		thresholds: [...{
			colorMode:  string
			op:         string
			$$hashKey?: string | null
			yaxis?:     string | null
			value:      number
			fill:       bool
			line:       bool
			visible?:   bool | null
		}]
		seriesOverrides: [...{
			alias:         string
			$$hashKey?:    string | null
			color?:        string | null
			dashLength?:   number | null
			fill?:         number | null
			linewidth?:    number | null
			yaxis?:        number | null
			zindex?:       number | null
			dashes?:       bool | null
			hiddenSeries?: bool | null
			hideTooltip?:  bool | null
		}]
		legend: {
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
		fieldConfig?: {
			overrides: []
			defaults: {
				unit?:      string | null
				unitScale?: bool | null
				links?: [] | null
				custom?: {} | null
			}
		} | null
		alert?: {
			executionErrorState: string
			for:                 string
			frequency:           string
			message:             string
			name:                string
			noDataState:         string
			handler:             number
			alertRuleTags: og_priority?: string | null
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
		title:            string
		repeat?:          string | null
		repeatDirection?: string | null
		id?:              number | null
		span?:            number | null
		collapsed:        bool
		editable?:        bool | null
		error?:           bool | null
		panels: [...#panel]
		targets?: [...#target] | null
		datasource?: #datasource
		gridPos?:    #grid
	}
	#_stat: {
		pluginVersion:  string
		title:          string
		description?:   string | null
		id?:            number | null
		maxDataPoints?: number | null
		transparent?:   bool | null
		targets: [...#target]
		links?: [] | null
		panels?: [...#panel] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			colorMode:          string
			graphMode:          string
			justifyMode:        string
			orientation:        string
			textMode:           string
			showPercentChange?: bool | null
			wideLayout?:        bool | null
			text?: {} | null
			reduceOptions: {
				fields: string
				values: bool
				calcs: [...string]
			}
		}
		transformations?: [...{
			id: string
			options: {
				indexByName: {}
				renameByName: {}
				excludeByName: {
					Time:  bool
					Value: bool
				}
			}
		}] | null
		fieldConfig: {
			overrides: [...{
				matcher: {
					id:      string
					options: string
				}
				properties: [...{
					id: string
					value: {
						mode: string
						steps: [...{
							color: string
							value: number | null
						}]
					}
				}]
			}]
			defaults: {
				unit?:      string | null
				min?:       number | null
				unitScale?: bool | null
				color?: {
					mode:        string
					fixedColor?: string | null
				} | null
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
				mappings: [...{
					type: string
					options: {
						match: string
						result: text: string
					}
				}]
			}
		}
	}
	#_table: {
		pluginVersion: string
		description?:  string | null
		title?:        string | null
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			cellHeight:  string
			frameIndex?: number | null
			showHeader:  bool
			sortBy?: [] | null
			footer: {
				fields:    string
				countRows: bool
				show:      bool
				reducer: [...string]
			}
		}
		transformations: [...{
			id: string
			options: {
				alias?:   string | null
				byField?: string | null
				match?:   string | null
				mode?:    string | null
				type?:    string | null
				fields?: {} | null
				includeByName?: {} | null
				reduce?: {
					reducer: string
				} | null
				binary?: {
					left:     string
					operator: string
					right:    string
				} | null
				excludeByName?: {
					Time?:  bool | null
					Value?: bool | null
				} | null
				include?: {
					names: [...string]
				} | null
				sort?: [...{
					field: string
				}] | null
				renameByName?: {
					"Value #A"?: string | null
					"Value #B"?: string | null
					"Value #C"?: string | null
					"Value #D"?: string | null
					"Value #E"?: string | null
					"Value #F"?: string | null
					"Value #G"?: string | null
					ip_address?: string | null
					task_name?:  string | null
				} | null
				indexByName?: {
					"Value #A"?:     number | null
					"Value #B"?:     number | null
					"Value #C"?:     number | null
					"Value #D"?:     number | null
					"Value #E"?:     number | null
					cluster_name?:   number | null
					monitoring_url?: number | null
				} | null
				filters?: [...{
					fieldName: string
					config: {
						id: string
						options: value: number
					}
				}] | null
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
						title?:       string | null
						type?:        string | null
						url?:         string | null
						targetBlank?: bool | null
						options?: {
							"1": {
								text:  string
								index: number
							}
							""?: {
								text:  string
								index: number
							} | null
						} | null
					}] | {
						mode?:             string | null
						type?:             string | null
						valueDisplayMode?: string | null
						steps?: [...{
							color:  string
							value?: number | null
						}] | null
					} | number | string
				}]
			}]
			defaults: {
				unit?:        string | null
				unitScale:    bool
				fieldMinMax?: bool | null
				color: mode: string
				custom: {
					align:       string
					inspect:     bool
					filterable?: bool | null
					cellOptions: type: string
				}
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
				mappings: [...{
					type: string
					options: {
						alive: {
							color: string
							text:  string
							index: number
						}
						probably_dead: {
							color: string
							text:  string
							index: number
						}
					}
				}]
			}
		}
	}
	#_text: {
		pluginVersion: string
		description?:  string | null
		title?:        string | null
		id?:           number | null
		span?:         number | null
		editable?:     bool | null
		error?:        bool | null
		transparent?:  bool | null
		links?: [] | null
		panels?: [...#panel] | null
		targets?: [...#target] | null
		style?: {
			"font-size": string
		} | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			content: string
			mode:    string
			code: {
				language:        string
				showLineNumbers: bool
				showMiniMap:     bool
			}
		}
	}
	#_timeseries: {
		title:            string
		description?:     string | null
		interval?:        string | null
		pluginVersion?:   string | null
		repeat?:          string | null
		repeatDirection?: string | null
		id?:              number | null
		maxPerRow?:       number | null
		targets: [...#target]
		panels?: [...#panel] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			tooltipOptions?: {
				mode: string
			} | null
			tooltip: {
				mode: string
				sort: string
			}
			legend: {
				displayMode: string
				placement:   string
				sortBy?:     string | null
				width?:      number | null
				showLegend:  bool
				sortDesc?:   bool | null
				calcs: [...string]
			}
		}
		transformations?: [...{
			id: string
			options: {
				match?:         string | null
				regex?:         string | null
				renamePattern?: string | null
				type?:          string | null
				filters?: [...{
					fieldName: string
					config: {
						id: string
						options: value: number
					}
				}] | null
			}
		}] | null
		alert?: {
			executionErrorState: string
			for:                 string
			frequency:           string
			message:             string
			name:                string
			noDataState:         string
			handler:             number
			alertRuleTags: og_priority?: string | null
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
		fieldConfig: {
			overrides: [...{
				"__systemRef"?: string | null
				matcher: {
					id: string
					options: {
						mode?:     string | null
						prefix?:   string | null
						readOnly?: bool | null
						names?: [...string] | null
					} | string
				}
				properties: [...{
					id: string
					value: {
						fill?:       string | null
						fixedColor?: string | null
						mode?:       string | null
						seriesBy?:   string | null
						graph?:      bool | null
						legend?:     bool | null
						tooltip?:    bool | null
						viz?:        bool | null
						dash?: [...number] | null
					} | number | string
				}]
			}]
			defaults: {
				unit?:      string | null
				decimals?:  number | null
				min?:       number | null
				unitScale?: bool | null
				links?: [] | null
				color: mode: string
				thresholds: {
					mode: string
					steps: [...{
						color:  string
						value?: number | null
					}]
				}
				mappings: [...{
					type: string
					options: {
						pattern: string
						result: {
							color: string
							index: number
						}
					}
				}]
				custom: {
					axisColorMode:     string
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
					axisCenteredZero:  bool
					spanNulls:         bool
					axisBorderShow?:   bool | null
					insertNulls?:      bool | null
					scaleDistribution: type: string
					thresholdsStyle: mode:   string
					stacking: {
						group: string
						mode:  string
					}
					hideFrom: {
						legend:  bool
						tooltip: bool
						viz:     bool
						graph?:  bool | null
					}
				}
			}
		}
	}
	type: "alertlist" | "bargauge" | "dashlist" | "gauge" | "graph" | "row" | "stat" | "table" | "text" | "timeseries"
	id?:  number
	panels: [...#panel]
	datasource?: #datasource
	gridPos?:    #grid
	targets: [...#target]
	if type == "alertlist" {
		#_alertlist
	}
	if type == "bargauge" {
		#_bargauge
	}
	if type == "dashlist" {
		#_dashlist
	}
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
