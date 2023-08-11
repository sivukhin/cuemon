package cuemon

#Conversion: {Input: #GrafanaPanel, Output: {
	if Input.type == "row" {
		Title:     Input.title
		Collapsed: Input.collapsed
	}
}}
