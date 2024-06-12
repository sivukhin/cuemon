package cuemon

import (
	listFunc "list"
	"strings"
	"strconv"
	"regexp"
)

#RowGriding: {
	Row: #Row
	SequencePanels: [for title, panel in Row.Panel {panel, Title: title}]
	SequenceGrid: [for i, panel in SequencePanels {
		Id:            i
		X:             (number & >=0 & <24) | *0
		Y:             (number & >=0) | *0
		Width:         number & >0 & <=24
		Height:        number & >0
		#EndX:         X + Width
		#ColumnNumber: number | *0
		#RowNumber:    number | *0

		if Row.Columns != _|_ && Row.Heights != _|_ {
			if Row.Columns[#ColumnNumber] != _|_ {
				Width: number | *Row.Columns[#ColumnNumber]
			}
			if Row.PanelGrid[panel.Title] != _|_ {
				if Row.PanelGrid[panel.Title].Height != _|_ {
					Height: Row.PanelGrid[panel.Title].Height
				}
				if Row.PanelGrid[panel.Title].Width != _|_ {
					Width: Row.PanelGrid[panel.Title].Width
				}
				if i > 0 {#ColumnNumber: -1}
			}
			if Row.PanelGrid[panel.Title] == _|_ {
				if i > 0 {
					if SequenceGrid[i-1].#ColumnNumber+1 < len(Row.Columns) {#ColumnNumber: SequenceGrid[i-1].#ColumnNumber + 1}
					if SequenceGrid[i-1].#ColumnNumber+1 == len(Row.Columns) {#ColumnNumber: 0}
				}
			}
		}

		if Row.Heights[#RowNumber] != _|_ {Height: number | *Row.Heights[#RowNumber]}
		if Row.Heights[#RowNumber] == _|_ && Row.Heights != _|_ {Height: number | *Row.Heights}
		if i > 0 {
			if SequenceGrid[i-1].#EndX+Width <= 24 {
				X:          SequenceGrid[i-1].#EndX
				Y:          SequenceGrid[i-1].Y
				#RowNumber: SequenceGrid[i-1].#RowNumber
			}
			if SequenceGrid[i-1].#EndX+Width > 24 {
				X:          0
				Y:          SequenceGrid[i-1].Y + SequenceGrid[i-1].Height
				#RowNumber: SequenceGrid[i-1].#RowNumber + 1
			}
		}
	}]
	RowGrided: {
		Title:     Row.Title
		Collapsed: Row.Collapsed
		Panel: {for i, position in SequenceGrid {
			"\(SequencePanels[i].Title)": {position, SequencePanels[i]}
		}}
		if len(SequenceGrid) > 0 {
			Height: listFunc.Max([for _, position in SequenceGrid {position.Y + position.Height}])
		}
		Count: len(SequenceGrid)
	}
}

#DashboardGriding: {
	DashboardRows: [...#Row]
	DashboardRowsGrided: [for row in Rows {{#RowGriding, Row: row}.RowGrided}]
	DashboardGrided: [for i, row in DashboardRowsGrided {row
		Id: number | *0
		Y:  number | *0
		if i > 0 {
			Id: DashboardGrided[i-1].#EndId
			Y:  DashboardGrided[i-1].#EndY
		}
		#EndId: Id + row.Count + 1
		#EndY:  number
		if row.Collapsed {
			#EndY: Y + 1
		}
		if !row.Collapsed {
			#EndY: Y + 1 + row.Height
		}
	}]
}

#Alphabet: [for x in listFunc.Range(10, 256+10, 1) {strings.ToUpper(strconv.FormatInt(x, 36))}]

#Panel2Grafana: {
	Row: {#Row, Y: number, Id: number}
	Panel: {#Panel, #Grid, Title: string, Id: number}
	G: {
		id:    Row.Id + 1 + Panel.Id
		title: Panel.Title
		type:  Panel.Type
		if Grafana.schemaVersion < 37 {datasource: Panel.DataSource}
		if Grafana.schemaVersion >= 37 {datasource: uid: Panel.DataSource}

		gridPos: {
			w: Panel.Width
			h: Panel.Height
			x: Panel.X
			if Row.Collapsed {
				y: Panel.Y
			}
			if !Row.Collapsed {
				y: Row.Y + 1 + Panel.Y
			}
		}
		if Panel.Alert != _|_ {
			let alertRegex = #"(?P<reducer>\w+)\((?P<ref>[^,]+),(?P<duration>1m|5m|10m|15m|1h),(?P<end>now|now-1m|now-5m)\) (?P<op>>|<|in|notin) [( ]*(?P<param1>[0-9.]+)[, ]*(?P<param2>[0-9.]+)?[) ]*"#
			let firstNotification = Panel.Alert.Notifications[0]
			let firstMatch = regexp.FindNamedSubmatch(alertRegex, firstNotification)
			if firstMatch.op == "in" {
				thresholds: [
					{value: strconv.ParseFloat(firstMatch.param1, 64), op: "gt"},
					{value: strconv.ParseFloat(firstMatch.param2, 64), op: "lt"},
				]
			}
			if firstMatch.op == "notin" {
				thresholds: [
					{value: strconv.ParseFloat(firstMatch.param1, 64), op: "lt"},
					{value: strconv.ParseFloat(firstMatch.param2, 64), op: "gt"},
				]
			}
			if firstMatch.op != "in" && firstMatch.op != "notin" {
				thresholds: [{
					value: strconv.ParseFloat(firstMatch.param1, 64)
					if firstMatch.op == ">" {op: "gt"}
					if firstMatch.op == "<" {op: "lt"}
				}]
			}

			alert: {
				alertRuleTags:       Panel.Alert.Tags
				executionErrorState: Panel.Alert.ExecutionErrorState
				noDataState:         Panel.Alert.NoDataState
				frequency:           Panel.Alert.Frequency
				"for":               Panel.Alert.PendingPeriod
				message:             Panel.Alert.Message
				name:                Panel.Alert.Name
				conditions: [for notification in Panel.Alert.Notifications {
					let match = regexp.FindNamedSubmatch(alertRegex, notification)
					if match.op == "in" {
						evaluator: params: [strconv.ParseFloat(match.param1, 64), strconv.ParseFloat(match.param2, 64)]
						evaluator: type: "within_range"
					}
					if match.op == "notin" {
						evaluator: params: [strconv.ParseFloat(match.param1, 64), strconv.ParseFloat(match.param2, 64)]
						evaluator: type: "outside_range"
					}
					if match.op == ">" {
						evaluator: params: [strconv.ParseFloat(match.param1, 64)]
						evaluator: type: "gt"
					}
					if match.op == "<" {
						evaluator: params: [strconv.ParseFloat(match.param1, 64)]
						evaluator: type: "lt"
					}
					reducer: type: match.reducer
					query: params: [match.ref, match.duration, match.end]
				}]

				if !Test {notifications: [for channel in Panel.Alert.Channels {uid: channel}]}
			}
		}
		if Panel.Type == "timeseries" {
			fieldConfig: defaults: unit: Panel.Unit
			let match = regexp.FindNamedSubmatch(#"(?P<displayMode>\w+)_(?P<placement>\w+)"#, Panel.Legend)
			options: legend: calcs:       Panel.Values
			options: legend: displayMode: match.displayMode
			options: legend: placement:   match.placement
		}
		if Panel.Type == "graph" {
			if Panel.Points > 0 {
				pointradius: Panel.Points
				points:      true
			}
			if Panel.Points == 0 {
				points: false
			}
			if Panel.Lines > 0 {
				linewidth: Panel.Lines
				lines:     true
			}
			if Panel.Lines == 0 {
				lines: false
			}
			stack:         Panel.Stack
			fill:          Panel.Fill
			nullPointMode: Panel.NullValue
			legend: {
				show:         Panel.Legend != "none"
				alignAsTable: Panel.Legend == "table_right" || Panel.Legend == "table_bottom"
				rightSide:    Panel.Legend == "table_right" || Panel.Legend == "list_right"
				values:       len(Panel.Values) > 0
				current:      listFunc.Contains(Panel.Values, "current")
				avg:          listFunc.Contains(Panel.Values, "avg")
				max:          listFunc.Contains(Panel.Values, "max")
				min:          listFunc.Contains(Panel.Values, "min")
				total:        listFunc.Contains(Panel.Values, "total")
				sort:         Panel.Sort
				sortDesc:     Panel.SortDesc
			}
			yaxes: [
				{$$hashKey: "object:\(10*(Row.Id+Panel.Id))", format:   Panel.Unit, if Panel.YMin != null {min: "\(Panel.YMin)"}},
				{$$hashKey: "object:\(10*(Row.Id+Panel.Id)+1)", format: Panel.Unit, if Panel.YMin != null {min: "\(Panel.YMin)"}},
			]
		}
		if Panel.Type == "stat" || Panel.Type == "gauge" {
			options: {
				textMode:  Panel.TextMode
				graphMode: Panel.GraphMode
				if Panel.Reduce == "all" {
					reduceOptions: values: true
				}
				if Panel.Reduce != "all" {
					reduceOptions: calcs: [Panel.Reduce]
					reduceOptions: values: false
				}
			}
		}
		if Panel.Thresholds != _|_ {
			fieldConfig: defaults: thresholds: steps: [for t in Panel.Thresholds {color: t.Color, value: t.Value}]
		}
		seriesOverrides: [for i, target in Panel.Metrics if target.Overrides != _|_ {
			$$hashKey: "object:\(10*(Row.Id+Panel.Id+2+i))"
			if target.Overrides.Alias != _|_ {alias: target.Overrides.Alias}
			if target.Overrides.Alias == _|_ {
				if strings.Contains(target.Legend, "{{") && strings.Contains(target.Legend, "}}") {
					alias: "/" + regexp.ReplaceAll("{{.*?}}", target.Legend, ".*") + "/"
				}
				if !strings.Contains(target.Legend, "{{") || !strings.Contains(target.Legend, "}}") {
					alias: target.Legend
				}
			}
			if target.Overrides.Dashes != _|_ {dashes: target.Overrides.Dashes}
			if target.Overrides.Hidden != _|_ {hiddenSeries: target.Overrides.Hidden}
			if target.Overrides.Fill != _|_ {fill: target.Overrides.Fill}
			if target.Overrides.YAxis != _|_ {yaxis: target.Overrides.YAxis}
			if target.Overrides.ZIndex != _|_ {zindex: target.Overrides.ZIndex}
			if target.Overrides.LineWidth != _|_ {linewidth: target.Overrides.LineWidth}
			if target.Overrides.Color != _|_ {color: target.Overrides.Color}
		}]
		targets: [for i, target in Panel.Metrics {
			refId: #Alphabet[i]
			hide:  target.Hide
			if target.StackDriver != _|_ {
				queryType: "metrics"
				metricQuery: {
					query:              target.Expr
					aliasBy:            target.Legend
					crossSeriesReducer: target.StackDriver.Reducer
					filters:            target.StackDriver.Filters
					groupBys:           target.StackDriver.GroupBy
					perSeriesAligner:   target.StackDriver.Aligner
					alignmentPeriod:    target.StackDriver.AlignmentPeriod
					projectName:        target.StackDriver.Project
					unit:               target.StackDriver.Unit
					valueType:          target.StackDriver.Value
					metricKind:         target.StackDriver.MetricKind
					metricType:         target.StackDriver.MetricType
					editorMode:         target.StackDriver.EditorMode
					if target.StackDriver.Preprocessor != _|_ {preprocessor: target.StackDriver.Preprocessor}
				}
			}
			if target.StackDriver == _|_ {
				expr:         target.Expr
				legendFormat: target.Legend
			}
		}]
	}
}

#variable2grafana: {
	name:     string
	variable: #variable
	grafana: {
		if variable.type == "constant" {
			type:  "constant"
			name:  name
			query: variable.Value
		}
		if variable.type == "custom" {
			type:       "custom"
			label:      variable.label
			name:       name
			query:      strings.Join(variable.values, ",")
			includeAll: variable.includeAll
			multi:      variable.multi
			current: {
				text:  variable.current
				value: variable.current
			}
			options: [for value in variable.values {
				selected: {
					if multi {listFunc.Contains(variable.current, value)}
					if !multi {variable.current == value}
				}
				text:  value
				value: value
			}]
		}
		if variable.type == "query" {
			type:       "query"
			label:      variable.label
			name:       name
			datasource: variable.dataSource
			definition: variable.query
			includeAll: variable.includeAll
			if variable.allValue != _|_ {allValue: variable.allValue}
			multi: variable.multi
			current: {
				text:  variable.current
				value: variable.current
			}
			sort: variable.sort
		}
	}
}

grafanaDashboard: #grafanaDashboard & {
	links: [for title, link in links {{type: "link", title: title, url: link}}]
	templating: list: [for name, variable in variables {
		{#variable2grafana, name: variableName, variable: variable}.grafana
	}]
	if !#playground {tags: tags}

	panels: listFunc.FlattenN([for row in {#DashboardGriding, DashboardRows: Rows}.DashboardGrided {
		if row.Collapsed {
			[{
				type:      "row"
				title:     row.Title
				id:        row.Id
				collapsed: true
				gridPos: {w: 24, h: 1, x: 0, y: row.Y}
				panels: [for panel in row.Panel {
					{#Panel2Grafana, Row: row, Panel: panel}.G
				}]
			}]
		}
		if (!row.Collapsed) {
			listFunc.Concat([[{
				type:      "row"
				title:     row.Title
				id:        row.Id
				collapsed: false
				gridPos: {w: 24, h: 1, x: 0, y: row.Y}
				panels: []
			}], [for panel in row.Panel {
				{#Panel2Grafana, Row: row, Panel: panel}.G
			}]])
		}
	}], 1)
}
