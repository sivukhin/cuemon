package cuemon

import (
	"list"
)

#Conversion: {Input: #GrafanaPanel, Output: {
	Type: Input.type
	if Input.alert != _|_ {
		Alert: {
			Tags:                Input.alert.alertRuleTags
			ExecutionErrorState: Input.alert.executionErrorState
			NoDataState:         Input.alert.noDataState
			Frequency:           Input.alert.frequency
			PendingPeriod:       Input.alert."for"
			Message:             Input.alert.message
			Name:                Input.alert.name
			Notifications: [ for condition in Input.alert.conditions {{
				#Operator: string
				if condition.evaluator.type == "gt" {#Operator: ">"}
				if condition.evaluator.type == "lt" {#Operator: "<"}

				N: "\(condition.reducer.type)(\(condition.query.params[0]),\(condition.query.params[1]),\(condition.query.params[2])) \(#Operator) \(condition.evaluator.params[0])"
			}.N}]
			Channels: [for notification in Input.alert.notifications { notification.uid }]
		}
	}
	if Type == "graph" {
		Unit: Input.yaxes[0].format
		if Input.points {
			Points: Input.pointradius
		}
		if !Input.points {
			Points: 0
		}
		if Input.lines {
			Lines: Input.linewidth
		}
		if !Input.lines {
			Lines: 0
		}
		NullValue: Input.nullPointMode
		if !Input.legend.show {
			Legend: "none"
		}
		if Input.legend.show && Input.legend.alignAsTable && Input.legend.rightSide {
			Legend: "table_right"
		}
		if Input.legend.show && Input.legend.alignAsTable && !Input.legend.rightSide {
			Legend: "table_bottom"
		}
		if Input.legend.show && !Input.legend.alignAsTable && Input.legend.rightSide {
			Legend: "list_right"
		}
		if Input.legend.show && !Input.legend.alignAsTable && !Input.legend.rightSide {
			Legend: "list_bottom"
		}
		Values: [ for key, value in Input.legend if list.Contains(#LegendValues, key) && value {key}]
	}
	if Input.type == "stat" || Input.type == "gauge" {
		TextMode:  Input.options.textMode
		GraphMode: Input.options.graphMode
		if Input.options.reduceOptions.values {
			Reduce: "all"
		}
		if !Input.options.reduceOptions.values {
			Reduce: Input.options.reduceOptions.calcs[0]
		}
		if Input.fieldConfig.defaults.thresholds.steps != _|_ {
			Thresholds: [ for t in Input.fieldConfig.defaults.thresholds.steps {Color: t.color, Value: t.value}]
		}
	}
	DataSource: Input.datasource
}}
