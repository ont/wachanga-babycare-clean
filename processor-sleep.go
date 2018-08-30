package main

import (
	"encoding/json"
	"log"
	"time"
)

type ProcessorSleep struct {
	*ProcessorBasic
}

func (p *ProcessorSleep) CanProcess(rec *Record) bool { return rec.EventType == "sleep" }

func (p *ProcessorSleep) Process(rec *Record) {
	var data struct {
		Reports []Report `json:"reports"`
	}
	err := json.Unmarshal(rec.RawValue, &data)
	if err != nil {
		log.Fatal(err)
		//p.api.Delete(rec)
	}

	if len(data.Reports) > 20 {
		p.api.Delete(rec)
	}

	if dur, ok := p.getDuration(data.Reports); ok {
		rec.JsonValue["duration"] = dur.Seconds()
		p.api.Save(rec)
	} else {
		p.api.Delete(rec)
	}
}

func (p *ProcessorSleep) getDuration(reports []Report) (time.Duration, bool) {
	var (
		state string
		stamp time.Time
		total time.Duration
	)

	for _, report := range reports {
		state = report.State
		if state == "asleep" {
			stamp = report.CreatedAt
		}

		if state == "awake" {
			dur := report.CreatedAt.Sub(stamp)
			if dur > 20*time.Hour {
				return 0, false
			}
			total += dur
		}
	}

	if state == "asleep" {
		return 0, false
	}

	return total, true
}
