package cuemon

import (
	"list"
)

#Conversion: {Input: #GrafanaGraph, Output: {
	Type: Input.type
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
		Values: [for key, value in Input.legend if list.Contains(#LegendValues, key) && value { key }]
	}
	if Input.type == "stat" || Input.type == "gauge" {
		TextMode: Input.options.textMode
		GraphMode: Input.options.graphMode
		if Input.options.reduceOptions.values {
			Reduce: "all"
		}
		if !Input.options.reduceOptions.values {
			Reduce: Input.options.reduceOptions.calcs[0]
		}
		if Input.fieldConfig.defaults.thresholds.steps != _|_ {
			Thresholds: [for t in Input.fieldConfig.defaults.thresholds.steps {Color: t.color, Value: t.value}]
		}
	}
	DataSource: Input.datasource
}}
