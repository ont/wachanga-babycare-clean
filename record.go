package main

import (
	"encoding/json"
	"log"
)

type Record struct {
	Id        string
	EventType string `db:"event_type"`
	RawValue  []byte `db:"raw_value"`

	JsonValue map[string]interface{} // parsed raw_value
}

func (r *Record) Parse() {
	err := json.Unmarshal(r.RawValue, &r.JsonValue)
	if err != nil {
		log.Fatalln(err)
	}
}

func (r *Record) Serialize() {
	data, err := json.Marshal(r.JsonValue)
	if err != nil {
		log.Fatalln(err)
	}
	r.RawValue = data
}

func (r *Record) GetString(key string) string {
	if ival, found := r.JsonValue[key]; found {
		val, _ := ival.(string)
		return val
	}
	return ""
}

func (r *Record) GetFloat(key string) float64 {
	if ival, found := r.JsonValue[key]; found {
		val, _ := ival.(float64)
		return val
	}
	return 0
}
