package cuemon

#conversion: {
	#hideTypes: ["none", "label", "variable"]
	#sortTypes: ["none", "alphAsc", "alphDesc", "numAsc", "numDesc", "ialphAsc", "ialphDesc"]
	#refreshTypes: ["none", "dashboardLoad", "timeRangeChange"]
	input: #template
	output: {
		input
		if input.type == "constant" {
			"#constant": {
				name:  input.name
				value: input.query
			}
		}
		if input.type == "custom" {
			"#custom": {
				name:    input.name
				label:   input.label
				multi:   input.multi
				all:     input.includeAll
				options: input.options
				hide:    #hideTypes[input.hide]
			}
		}
		if input.type == "query" {
			"#datasrc": input.datasource
			"#query": {
				name:       input.name
				if (input.label & string) != _|_ {
					label: input.label
				}
				hide:       #hideTypes[input.hide]
				all:        input.includeAll
				multi:      input.multi
				sort:       #sortTypes[input.sort]
				regex:      input.regex
				refresh:    #refreshTypes[input.refresh]
				current:    input.current.value
				query:      input.query.query
				if (input.allValue & string) != _|_ {
					allValue: input.allValue
				}
			}
		}
	}
}
