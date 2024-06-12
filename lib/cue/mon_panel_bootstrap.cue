package cuemon

#conversion: {
	input: #panel
	output: {
		input
		if input.type == "text" {
			"#text": {
				title:       input.title
				markdown:    input.options.content
				transparent: input.transparent
			}
		}

		if input.type == "stat" {
			"#stat": {
				title:   input.title
				reducer: input.options.reduceOptions.calcs
			}
		}

		if input.type == "graph" {
			"#graph": {
				if input.legend != _|_ {
					if input.legend.alignAsTable != _|_ {
						if input.legend.alignAsTable == true {legend: type: "table"}
						if input.legend.alignAsTable == false {legend: type: "list"}
					}
					if input.legend.rightSide != _|_ {
						if input.legend.rightSide == true {legend: placement: "right"}
						if input.legend.rightSide == false {legend: placement: "bottom"}
					}
					if input.legend.sortDesc != _|_ {
						if input.legend.sortDesc == true {legend: sortHow: "desc"}
						if input.legend.sortDesc == false {legend: sortHow: "asc"}
					}
					if input.legend.sort != _|_ {
						legend: sortBy: input.legend.sort
					}
					legend: values: [for value in ["avg", "max", "min", "current", "total"] if input.legend[value] {value}]
				}
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
			targets: [for target in input.targets {
				if target.expr != _|_ {"#prom": target.expr}
				if target.legendFormat != _|_ {"#format": target.legendFormat}
				if target.metricQuery != _|_ {
					if target.metricQuery.editorMode == "mql" {
						"#mqlScript": {
							aliasBy:            target.metricQuery.aliasBy
							alignmentPeriod:    target.metricQuery.alignmentPeriod
							crossSeriesReducer: target.metricQuery.crossSeriesReducer
							perSeriesAligner:   target.metricQuery.perSeriesAligner
							projectName:        target.metricQuery.projectName
							query:              target.metricQuery.query
						}
					}
				}
			}]
		}
	}
}
