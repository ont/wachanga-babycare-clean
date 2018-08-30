package main

type Processor interface {
	CanProcess(rec *Record) bool
	Process(rec *Record)
}

type ProcessorBasic struct {
	api *DBApi
}

func NewProcessorBasic(api *DBApi) *ProcessorBasic {
	return &ProcessorBasic{api}
}
