package cuemon

#Conversion: {SchemaVersion: number, Input: #GrafanaTemplate, Output: {
	Type: Input.type
	if Type == "constant" {
		Value: Input.query
	}
	if Type == "custom" {
		Label: Input.label
		Values: [ for option in Input.options if option.value != "$__all" {option.value}]
		Multi:      Input.multi
		IncludeAll: Input.includeAll
		Current:    Input.current.value
	}
	if Type == "query" {
		Label: Input.label
		DataSource: Input.datasource
		Query:      Input.definition
		Multi:      Input.multi
		IncludeAll: Input.includeAll
		Current:    Input.current.value
		Sort:       Input.sort
	}
}}
