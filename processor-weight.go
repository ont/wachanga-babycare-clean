package main

type ProcessorWeight struct {
	*ProcessorBasic
}

func (p *ProcessorWeight) CanProcess(rec *Record) bool {
	return rec.EventType == "measurement" && rec.GetString("measurement_type") == "weight"
}

func (p *ProcessorWeight) Process(rec *Record) {
	value := rec.GetFloat("value")

	if value > 30000 {
		value = value / 10
		rec.JsonValue["value"] = value
		p.api.Save(rec)
	}

	if value < 1000 {
		p.api.Delete(rec)
	}
}
