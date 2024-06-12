package cuemon

import (
	"list"
)

#gridPanel: {
	#w?: number
	#h?: number
	...
}

#gridGroup: {
	#height: number

	#columns: [...number] & list.MinItems(1) // cyclic array of column widths (so, with [8] - all columns will have width of 8 units)

	panels: [...#gridPanel]
}

#gridRow: {
	rowHeight=#height?: number
	rowColumns=#columns?: [...number] & list.MinItems(1)

	collapsed: bool | *false
	groups: [...(#gridGroup & {#height: number | *rowHeight, #columns: [...number] | *rowColumns, ...})]
}

#monGrid: {
	rows: [...#gridRow]
}
