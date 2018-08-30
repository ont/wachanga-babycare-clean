package main

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
	"time"
)

type ProcessorLactation struct {
	*ProcessorBasic
}

func (p *ProcessorLactation) CanProcess(rec *Record) bool { return rec.EventType == "lactation" }

func (p *ProcessorLactation) Process(rec *Record) {
	var data struct {
		Reports []Report `json:"reports"`
	}
	err := json.Unmarshal(rec.RawValue, &data)
	if err != nil {
		log.Fatal(err)
		// p.api.Delete(rec)
	}

	if len(data.Reports) > 20 {
		p.api.Delete(rec)
	}

	if durLeft, durRight, ok := p.getDurations(data.Reports); ok {
		//spew.Dump(durLeft+durRight, durLeft, durRight, data.Reports)
		// save duration to rec
		rec.JsonValue["duration_left"] = durLeft.Seconds()
		rec.JsonValue["duration_right"] = durRight.Seconds()
		rec.JsonValue["duration"] = (durRight + durLeft).Seconds()
		p.api.Save(rec)
	} else {
		p.api.Delete(rec)
	}
}

func (p *ProcessorLactation) getDurations(reports []Report) (time.Duration, time.Duration, bool) {
	sort.Sort(ReportsByCreated(reports))

	state := map[string]string{"left": "stop", "right": "stop"}
	stamps := map[string]time.Time{}
	durs := map[string]time.Duration{"left": 0, "right": 0}

	for _, report := range reports {
		parts := strings.Split(report.State, "_")

		//if parts[1] == state[parts[0]] {
		//	fmt.Println("wrong state")
		//	return 0, 0, false
		//}

		state[parts[0]] = parts[1]

		if state[parts[0]] == "start" {
			stamps[parts[0]] = report.CreatedAt
		}

		if state[parts[0]] == "stop" {
			dur := report.CreatedAt.Sub(stamps[parts[0]])
			if dur > 30*time.Minute {
				return 0, 0, false
			}
			durs[parts[0]] += dur
		}
	}

	if state["left"] == "start" || state["right"] == "start" {
		//fmt.Println("not all stopped")
		return 0, 0, false
	}

	return durs["left"], durs["right"], true
}
