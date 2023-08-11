package cuemon

#Conversion: {Input: #GrafanaSchema, Output: {
	Grafana: id:  Input.id
	Grafana: uid: Input.uid

	Title: Input.title

	Links: {for link in Input.links if link.type == "link" {
		"\(link.title)": Url: link.url
	}}
	Tags: Input.tags
}}
