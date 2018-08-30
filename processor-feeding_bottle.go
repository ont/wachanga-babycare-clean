package main

type ProcessorFeedingBottle struct {
	*ProcessorBasic
}

func (p *ProcessorFeedingBottle) CanProcess(rec *Record) bool {
	return rec.EventType == "feeding_bottle"
}

func (p *ProcessorFeedingBottle) Process(rec *Record) {
	volume := rec.GetFloat("volume")

	if volume > 1000 {
		rec.JsonValue["volume"] = volume / 10
		p.api.Save(rec)
	}
}
