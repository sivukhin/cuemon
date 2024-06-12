package cuemon

import (
	"list"
)

#conversion: {
	#filter: ["templating", "annotations", "links", "panels"]
	input: #grafana
	output: {
		"#meta": {
			id:       input.id
			uid:      input.uid
			title:    input.title
			refresh:  input.refresh
			from:     input.time.from
			to:       input.time.to
			timezone: input.timezone
			tags:     input.tags
			if input.timepicker.refresh_intervals != _|_ {timepicker: input.timepicker.refresh_intervals}
			if #grafanaVersion == "v10" && input.liveNow != _|_ {live: input.liveNow}
		}
		for key, value in input if !list.Contains(#filter, key) {
			"\(key)": value
		}
	}
}
