package cuemon

#Conversion: {Grafana:#GrafanaGraph, Mon: {
	Type: Grafana.type
	if Type == "graph" {
		Unit:       Grafana.yaxes[0].format
	}
	DataSource: Grafana.datasource
	Metrics: [for target in Grafana.targets {
		if target.queryType == _|_ {
			Expr: target.expr
			Legend: target.legendFormat
		}
		if target.queryType != _|_ {
			if target.queryType == "randomWalk" {
				Expr: target.expr
				Legend: target.legendFormat
			}
			if target.queryType == "metrics" {
				Expr: target.metricQuery.query
				Legend: target.metricQuery.aliasBy
				StackDriver: {
					Reducer: target.metricQuery.crossSeriesReducer
					if len(target.metricQuery.filters) > 0 {
						Filters: target.metricQuery.filters
					}
					if len(target.metricQuery.groupBys) > 0 {
						GroupBy: target.metricQuery.groupBys
					}
					Aligner: target.metricQuery.perSeriesAligner
					Project: target.metricQuery.projectName
					if target.metricQuery.unit != "" {
						Uint: target.metricQuery.unit
					}
					if target.metricQuery.valueType != "" {
						Value: target.metricQuery.valueType
					}
				}
			}
		}
	}]
}}
