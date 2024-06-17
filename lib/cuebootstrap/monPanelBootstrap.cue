package cuemon

import (
	"strconv"
)

#conversion: {
	#v7ToV10: {min: "min", max: "max", mean: "mean", avg: "mean", total: "sum", sum: "sum", current: "lastNotNull", lastNotNull: "lastNotNull", last: "last"}
	input: #panel
	output: {
		input
		"#title": input.title
		if input.description != _|_ {"#description": input.description}
		if input.datasource != _|_ {"#datasrc": input.datasource}

		if input.type == "text" {
			"#text": {
				markdown:    input.options.content
				transparent: input.transparent
			}
		}

		if input.type == "stat" {
			"#stat": {
				reducer: input.options.reduceOptions.calcs
				if input.fieldConfig.defaults.unit != _|_ {yPrimary: unit: input.fieldConfig.defaults.unit}
				if input.fieldConfig.defaults.min != _|_ {yPrimary: min: input.fieldConfig.defaults.min}
				if input.fieldConfig.defaults.max != _|_ {yPrimary: min: input.fieldConfig.defaults.max}
			}
		}

		if input.alert != _|_ {
			"#alert": input.alert
		}

		if input.seriesOverrides != _|_ {
			"#overrides": [for override in input.seriesOverrides {
				hashKey: override.$$hashKey
				alias:   override.alias
				if override.color != _|_ {color: override.color}
				if override.linewidth != _|_ {linewidth: override.linewidth}
				if override.yaxis != _|_ {yaxis: override.yaxis}
				if override.hiddenSeries != _|_ {hidden: override.hiddenSeries}
				if override.dashes != _|_ {dashes: override.dashes}
				if override.legend != _|_ {legend: override.legend}
				if override.fillGradient != _|_ {fillGradient: override.fillGradient}
			}]
		}

		if input.type == "timeseries" {
			if #grafanaVersion == "v7" {"#defaultGraphPlugin": false}
			"#graph": {
				legend: type:      input.options.legend.displayMode
				legend: placement: input.options.legend.placement
				legend: values:    input.options.legend.calcs
			}
			//			if input.fieldConfig.defaults.min != _|_ {"#rangeY": min: input.fieldConfig.defaults.min}
			//			if input.fieldConfig.defaults.max != _|_ {"#rangeY": min: input.fieldConfig.defaults.max}
		}

		if input.type == "graph" {
			if #grafanaVersion == "v10" {"#defaultGraphPlugin": false}
			"#graph": {
				if input.legend != _|_ {
					if input.legend.alignAsTable != _|_ {
						if input.legend.alignAsTable == true {legend: type: "table"}
						if input.legend.alignAsTable == false {legend: type: "list"}
					}
					if input.legend.rightSide != _|_ {
						if input.legend.rightSide == true {legend: placement: "right"}
						if input.legend.rightSide != true {legend: placement: "bottom"}
					}
					if input.legend.rightSide == _|_ {legend: placement: "bottom"}
					if input.legend.sortDesc != _|_ {
						if input.legend.sortDesc == true {legend: sortHow: "desc"}
						if input.legend.sortDesc == false {legend: sortHow: "asc"}
					}
					if input.legend.sort != _|_ && (input.legend.sort & string) != _|_ {
						legend: sortBy: #v7ToV10[input.legend.sort]
					}
					legend: values: [for value in ["avg", "max", "min", "current", "total"] if input.legend[value] {#v7ToV10[value]}]
					legend: showValues: input.legend.values
				}
				yPrimary: unit:      input.yaxes[0].format
				yPrimary: hashKey:   input.yaxes[0].$$hashKey
				ySecondary: unit:    input.yaxes[1].format
				ySecondary: hashKey: input.yaxes[1].$$hashKey
				if (input.yaxes[0].min & string) != _|_ {yPrimary: min: strconv.ParseFloat(input.yaxes[0].min, 64)}
				if (input.yaxes[0].max & string) != _|_ {yPrimary: max: strconv.ParseFloat(input.yaxes[0].max, 64)}
				if (input.yaxes[1].min & string) != _|_ {ySecondary: min: strconv.ParseFloat(input.yaxes[1].min, 64)}
				if (input.yaxes[1].max & string) != _|_ {ySecondary: max: strconv.ParseFloat(input.yaxes[1].max, 64)}
				display: fill: input.fill
				display: bars: use:    input.bars
				display: lines: use:   input.lines
				display: lines: size:  input.linewidth
				display: points: use:  input.points
				display: points: size: input.pointradius
				display: stack: input.stack
				if input.nullPointMode == "connected" {display: nulls: "connected"}
				if input.tooltip != _|_ {tooltip: sort: input.tooltip.sort}
			}
		}

		if input.fieldConfig.defaults.unit != _|_ {
			"#unit": input.fieldConfig.defaults.unit
		}

		if input.fieldConfig.defaults.thresholds != _|_ {
			"#thresholds": {
				mode:  input.fieldConfig.defaults.thresholds.mode
				steps: input.fieldConfig.defaults.thresholds.steps
			}
		}

		if input.targets != _|_ {
			"#targets": [for target in input.targets {
				if target.hide != _|_ {hide: target.hide}
				if target.expr != _|_ {prom: target.expr}
				if target.legendFormat != _|_ {format: target.legendFormat}
				if target.metricQuery != _|_ {
					if target.metricQuery.editorMode == "mql" {
						mqlScript: {
							aliasBy:            target.metricQuery.aliasBy
							alignmentPeriod:    target.metricQuery.alignmentPeriod
							crossSeriesReducer: target.metricQuery.crossSeriesReducer
							perSeriesAligner:   target.metricQuery.perSeriesAligner
							projectName:        target.metricQuery.projectName
							query:              target.metricQuery.query
						}
					}
					if target.metricQuery.editorMode == "visual" {
						mqlVisual: {
							aliasBy:            target.metricQuery.aliasBy
							alignmentPeriod:    target.metricQuery.alignmentPeriod
							crossSeriesReducer: target.metricQuery.crossSeriesReducer
							filters:            target.metricQuery.filters
							groupBys:           target.metricQuery.groupBys
							metricKind:         target.metricQuery.metricKind
							metricType:         target.metricQuery.metricType
							perSeriesAligner:   target.metricQuery.perSeriesAligner
							projectName:        target.metricQuery.projectName
							unit:               target.metricQuery.unit
							valueType:          target.metricQuery.valueType
							if #grafanaVersion == "v7" {sloV7: target.sloQuery != _|_}
						}
					}
				}
			}]
		}
	}
}
