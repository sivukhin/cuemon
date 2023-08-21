package cuemon

#Conversion: {SchemaVersion: number, Input: #GrafanaSchema, Output: {
	Grafana: schemaVersion: Input.schemaVersion
	Grafana: title: Input.title
	Grafana: time: from: Input.time.from
	Grafana: time: to: Input.time.to

	Links: {for link in Input.links if link.type == "link" {"\(link.title)": Url: link.url}}
	Tags: Input.tags
}}
