package cuemon

#Conversion: {Grafana: #GrafanaTemplate, Mon: {
	Type: Grafana.type
	if Type == "constant" {
		Value: Grafana.query
	}
	if Type == "custom" {
		Values: [for option in Grafana.options { option.value }]
		Multi: Grafana.multi
		IncludeAll: Grafana.includeAll
		Current: Grafana.current.value
	}
	if Type == "query" {
		DataSource: Grafana.datasource
		Query: Grafana.definition
		Multi: Grafana.multi
		IncludeAll: Grafana.includeAll
		Current: Grafana.current.value
		Sort: Grafana.sort
	}
}}
