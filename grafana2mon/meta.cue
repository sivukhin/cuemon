package cuemon

#Conversion: {G=Grafana: #GrafanaSchema, Mon: {
	Grafana: id: G.id
	Grafana: uid: G.uid

	Title: G.title

	Links: { for link in G.links if link.type == "link" {
		"\(link.title)": Url: link.url
	}}
	Tags: G.tags
}}
