mon: #variables: [{
	#constant: {
		name:  "cluster1"
		value: "moj-p-ds-generic-services-02"
	}
}, {
	#query: {
		name:       "cluster2"
		datasource: #ds
		query:      "label_values(feature_service_error_rate_counter{}, service_prometheus_track)"
		current: ["/searchVectorNN"]
	}
}, {
	#textbox: {
		name:    "cluster3"
		default: "kek.*"
	}
}, {
	#custom: {
		name: "cluster4"
		options: [{value: "1", selected: true}, {value: "2"}]
	}
}]
