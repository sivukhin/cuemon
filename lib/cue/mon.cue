package cuemon

#mon: #grafana & {#grafana}
#mon: {
	#monMeta
	#links: [...#monLink]
	#variables: [...#monVariable]
	#rowsRef: [...#monRow]

	links: #links
	templating: list: #variables
}
