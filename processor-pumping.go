package main

type ProcessorPumping struct {
	*ProcessorBasic
}

func (p *ProcessorPumping) CanProcess(rec *Record) bool { return rec.EventType == "pumping" }

func (p *ProcessorPumping) Process(rec *Record) {
	volume := rec.GetFloat("volume")

	if volume > 1000 {
		rec.JsonValue["volume"] = volume / 10
		p.api.Save(rec)
	}
}
