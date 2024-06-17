package cuemon

#playground: bool @tag(playground,type=bool)

#mon: #grafana & {#grafana}
#mon: {
	#monMeta
	#links: [...#monLink]
	#variables: [...#monVariable]
	#rows: [...#monRow]

	links: #links
	templating: list: #variables
}
