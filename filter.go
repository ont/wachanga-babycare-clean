package main

type Filter struct {
	procs []Processor
	ch    <-chan *Record
}

type Tester func(*Record) bool

func NewFilter(ch chan *Record, procs []Processor) *Filter {
	return &Filter{procs, ch}
}

func (f *Filter) Run() {
	for rec := range f.ch {
		for _, proc := range f.procs {
			if proc.CanProcess(rec) {
				proc.Process(rec)
				break
			}
		}
	}
}
