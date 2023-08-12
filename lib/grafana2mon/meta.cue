package cuemon

#Conversion: {SchemaVersion: number, Input: #GrafanaSchema, Output: {
	Grafana: id:            Input.id
	Grafana: uid:           Input.uid
	Grafana: schemaVersion: Input.schemaVersion

	Title: Input.title

	Links: {for link in Input.links if link.type == "link" {
		"\(link.title)": Url: link.url
	}}
	Tags: Input.tags
}}
