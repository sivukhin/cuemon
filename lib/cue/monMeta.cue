package cuemon

#monMeta: #grafana & {#grafana}
#monMeta: {
	// cue trim -s works poorly if #meta will be defined (it trim too much and this is known bug)
	// so, here we introduce potentially undefined field in order to fix trim behaviour
	#meta?: {
		id:       number
		uid:      string
		title:    string
		refresh:  string | *"15m"
		from:     string | *"now-6h"
		to:       string | *"now"
		timezone: string | *"browser"
		tags: [...string]
		live: bool | *false
		timepicker?: [...string]
	}
	if #meta != _|_ {
		id:            _ | *#meta.id
		uid:           _ | *#meta.uid
		title:         _ | *#meta.title
		refresh:       _ | *#meta.refresh
		schemaVersion: _ | *#schemaVersion
		timezone:      _ | *#meta.timezone
		tags:          _ | *#meta.tags
		time: from: _ | *#meta.from
		time: to:   _ | *#meta.to
	}
	editable:     number | *true
	graphTooltip: 0 | *1 | 2 // default | shared-crosshair | shared-tooltpi
	if #meta.timepicker != _|_ {timepicker: refresh_intervals: #meta.timepicker}
	if #grafanaVersion == "v10" {
		fiscalYearStartMonth: number | *0
		liveNow:              bool | *#meta.live
		weekStart:            string | *""
	}
}
