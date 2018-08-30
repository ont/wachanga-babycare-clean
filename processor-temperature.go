package main

type ProcessorTemperature struct {
	*ProcessorBasic
}

func (p *ProcessorTemperature) CanProcess(rec *Record) bool {
	return rec.EventType == "temperature"
}

func (p *ProcessorTemperature) Process(rec *Record) {
	value := rec.GetFloat("value")

	if 1 <= value && value <= 6 {
		value = value*9/5 + 32
		rec.JsonValue["value"] = value
		p.api.Save(rec)
	}

	// after possible correction check for outbounds
	if value < 32 || value > 42 {
		p.api.Delete(rec)
	}
}
