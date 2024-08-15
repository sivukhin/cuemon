#grafana: {
	#grid: _
	#annotation: _
	#datasource: _
	#target: datasource: #datasource
	#template: {
		type: string @discriminative()
		datasource: #datasource
	}
	#yaxes: $$hashKey?: string
	#link: type: string @discriminative()
	#panel: {
		type: string @discriminative()
		gridPos?: #grid
		id?: number
		datasource: #datasource
		targets: [...#target]
		panels: [...#panel]
		yaxes?: [...#yaxes]
	}

	iteration?: number
	version?: number

	annotations: list: [...#annotation]
	links: [...#link]
	panels: [...#panel]
	templating: list: [...#template]
} @root(undefined-is-null,null-is-undefined)
