package analyzer

import (
	"sort"
	"time"
)

type TimingEntry struct {
	Task string `json:"task"`
	Time string `json:"time"`
}

// ConvertTimings converts map of timings to a sorted slice (desc by duration)
func ConvertTimings(timings map[string]string) []TimingEntry {
	if timings == nil {
		return nil
	}

	type sortable struct {
		task string
		time string
		dur  float64
	}

	list := make([]sortable, 0, len(timings))
	for k, v := range timings {
		d, err := time.ParseDuration(v)
		if err != nil {
			continue
		}
		list = append(list, sortable{task: k, time: v, dur: d.Seconds()})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].dur > list[j].dur // descending
	})

	result := make([]TimingEntry, len(list))
	for i, t := range list {
		result[i] = TimingEntry{Task: t.task, Time: t.time}
	}
	return result
}
