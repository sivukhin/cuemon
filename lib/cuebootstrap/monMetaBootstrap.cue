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
			from:     input.time.from
			to:       input.time.to
			timezone: input.timezone
			tags:     input.tags
			if (input.refresh & string) != _|_ {refresh: input.refresh}
			if input.timepicker.refresh_intervals != _|_ {timepicker: input.timepicker.refresh_intervals}
			if input.liveNow != _|_ {live: input.liveNow}
		}
		for key, value in input if !list.Contains(#filter, key) {
			"\(key)": value
		}
	}
}
