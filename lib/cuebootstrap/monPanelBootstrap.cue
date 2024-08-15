package cuemon

#conversion: {
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

		if input.type == "timeseries" {
			"#timeseries": {
			}
			//			if input.fieldConfig.defaults.min != _|_ {"#rangeY": min: input.fieldConfig.defaults.min}
			//			if input.fieldConfig.defaults.max != _|_ {"#rangeY": min: input.fieldConfig.defaults.max}
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
						}
					}
				}
			}]
		}
	}
}
