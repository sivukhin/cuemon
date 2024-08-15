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
	timepicker: refresh_intervals?: [...string] | null
	links: [...#link]
	annotations: list: [...#annotation]
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
	refId:                string
	aliasBy?:             string | null
	editorMode?:          string | null
	expr?:                string | null
	expression?:          string | null
	format?:              string | null
	interval?:            string | null
	legendFormat?:        string | null
	queryType?:           string | null
	type?:                string | null
	disableTextWrap?:     bool | null
	exemplar?:            bool | null
	fullMetaSearch?:      bool | null
	hide?:                bool | null
	includeNullMetadata?: bool | null
	instant?:             bool | null
	range?:               bool | null
	useBackend?:          bool | null
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
	timeSeriesList?: {
		alignmentPeriod:    string
		crossSeriesReducer: string
		perSeriesAligner:   string
		preprocessor:       string
		projectName:        string
		filters: [...string]
		groupBys: [...string]
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
		valueType:          string
		expr?:              string | null
		legendFormat?:      string | null
		preprocessor?:      string | null
		unit?:              string | null
		range?:             bool | null
		filters?: [...string] | null
		groupBys?: [...string] | null
	} | null
}
#template: {
	#_adhoc: {
		name:        string
		hide:        number
		skipUrlSync: bool
		filters: []
		datasource?: #datasource
	}
	#_constant: {
		name:        string
		query:       string
		hide:        number
		skipUrlSync: bool
		datasource?: #datasource
	}
	#_custom: {
		name:        string
		query:       string
		allValue?:   string | null
		label?:      string | null
		queryValue?: string | null
		hide:        number
		includeAll:  bool
		multi:       bool
		skipUrlSync: bool
		datasource?: #datasource
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
		queryValue:  string
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
		allValue?:       string | null
		description?:    string | null
		label?:          string | null
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
	type:        "adhoc" | "constant" | "custom" | "datasource" | "interval" | "query" | "textbox"
	datasource?: #datasource
	if type == "adhoc" {
		#_adhoc
	}
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
#yaxes: {
	$$hashKey?: string
}
#link: {
	#_link: {
		icon:        string
		title:       string
		tooltip:     string
		url:         string
		asDropdown:  bool
		includeVars: bool
		keepTime:    bool
		targetBlank: bool
		tags: []
	}
	type: "link"
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
		yaxes?: [...#yaxes] | null
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
		yaxes?: [...#yaxes] | null
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
			text: {}
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
				max:       number
				min:       number
				unitScale: bool
				mappings: []
				color: mode: string
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
	#_dashlist: {
		pluginVersion: string
		title:         string
		id?:           number | null
		panels?: [...#panel] | null
		targets?: [...#target] | null
		yaxes?: [...#yaxes] | null
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
		yaxes?: [...#yaxes] | null
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
						color:  string
						value?: number | null
					}]
				}
			}
		}
	}
	#_row: {
		title:            string
		repeat?:          string | null
		repeatDirection?: string | null
		id?:              number | null
		collapsed:        bool
		panels?: [...#panel] | null
		targets?: [...#target] | null
		yaxes?: [...#yaxes] | null
		datasource?: #datasource
		gridPos?:    #grid
	}
	#_stat: {
		pluginVersion: string
		title:         string
		description?:  string | null
		interval?:     string | null
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		yaxes?: [...#yaxes] | null
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
				regex?:         string | null
				renamePattern?: string | null
				indexByName?: {} | null
				renameByName?: {} | null
				excludeByName?: {
					Time:  bool
					Value: bool
				} | null
			}
		}] | null
		fieldConfig: {
			defaults: {
				unit?:      string | null
				min?:       number | null
				unitScale?: bool | null
				mappings: []
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
							color:  string
							value?: number | null
						}]
					}
				}]
			}]
		}
	}
	#_table: {
		pluginVersion: string
		description?:  string | null
		title?:        string | null
		id?:           number | null
		targets: [...#target]
		panels?: [...#panel] | null
		yaxes?: [...#yaxes] | null
		datasource?: #datasource
		gridPos?:    #grid
		options: {
			cellHeight:  string
			frameIndex?: number | null
			showHeader:  bool
			sortBy?: [...{
				displayName: string
				desc:        bool
			}] | null
			footer: {
				fields:            string
				countRows:         bool
				show:              bool
				enablePagination?: bool | null
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
				include?: {
					names: [...string]
				} | null
				sort?: [...{
					field: string
				}] | null
				renameByName?: {
					Value?:              string | null
					"Value #A"?:         string | null
					"Value #B"?:         string | null
					"Value #C"?:         string | null
					"Value #D"?:         string | null
					"Value #E"?:         string | null
					"Value #F"?:         string | null
					"Value #G"?:         string | null
					ip_address?:         string | null
					redpanda_group?:     string | null
					redpanda_partition?: string | null
					redpanda_topic?:     string | null
					task_name?:          string | null
				} | null
				excludeByName?: {
					Time?:     bool | null
					Value?:    bool | null
					instance?: bool | null
					job?:      bool | null
					shard?:    bool | null
				} | null
				indexByName?: {
					Time?:               number | null
					Value?:              number | null
					"Value #A"?:         number | null
					"Value #B"?:         number | null
					"Value #C"?:         number | null
					"Value #D"?:         number | null
					"Value #E"?:         number | null
					cluster_name?:       number | null
					monitoring_url?:     number | null
					redpanda_group?:     number | null
					redpanda_partition?: number | null
					redpanda_topic?:     number | null
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
				min?:         number | null
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
		description:   string
		pluginVersion: string
		title:         string
		id?:           number | null
		panels?: [...#panel] | null
		targets?: [...#target] | null
		yaxes?: [...#yaxes] | null
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
		yaxes?: [...#yaxes] | null
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
					lineStyle?: {
						fill: string
					} | null
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
	type: "alertlist" | "bargauge" | "dashlist" | "gauge" | "row" | "stat" | "table" | "text" | "timeseries"
	id?:  number
	panels: [...#panel]
	datasource?: #datasource
	gridPos?:    #grid
	yaxes?: [...#yaxes]
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
