package cuemon

#ds: _
if #grafanaVersion == "v7" {#ds: "VM-Services"}
if #grafanaVersion == "v10" {#ds: {type: "prometheus", uid: "TjtJft04z"}}

mon: #mon
mon: #meta: {
	title: "sivukhin - test update"
	id:    384                                    // 676
	uid:   "fd0b1f77-876d-45a7-9aae-14f0aed494f4" // "r3vajs54z"
}

mon: #links: [{
	title: "Hamsa", url: "http://..."
}]

mon: #rows: [#row1]
