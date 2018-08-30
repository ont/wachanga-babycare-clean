package main

type ProcessorGrowth struct {
	*ProcessorBasic
}

func (p *ProcessorGrowth) CanProcess(rec *Record) bool {
	return rec.EventType == "measurement" && rec.GetString("measurement_type") == "growth"
}

func (p *ProcessorGrowth) Process(rec *Record) {
	value := rec.GetFloat("value")
	if value > 400 {
		value = value / 10
		rec.JsonValue["value"] = value
		p.api.Save(rec)
	}

	if value < 20 || value > 130 {
		p.api.Delete(rec)
	}
}
