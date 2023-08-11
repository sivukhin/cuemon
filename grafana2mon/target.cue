package cuemon

#Conversion: {Input: #GrafanaTarget, Output: {
	Expr:   Input.expr | Input.metricQuery.query
	Legend: Input.legendFormat | Input.metricQuery.aliasBy
	if Input.queryType != _|_ {
		if Input.queryType == "metrics" {
			Expr:   Input.metricQuery.query
			Legend: Input.metricQuery.aliasBy
			StackDriver: {
				Reducer:         Input.metricQuery.crossSeriesReducer
				Filters:         Input.metricQuery.filters
				GroupBy:         Input.metricQuery.groupBys
				Aligner:         Input.metricQuery.perSeriesAligner
				AlignmentPeriod: Input.metricQuery.alignmentPeriod
				MetricKind:      Input.metricQuery.metricKind
				MetricType:      Input.metricQuery.metricType
				EditorMode:      Input.metricQuery.editorMode
				Project:         Input.metricQuery.projectName
				Unit:            Input.metricQuery.unit
				Value:           Input.metricQuery.valueType
			}
		}
	}
}}
