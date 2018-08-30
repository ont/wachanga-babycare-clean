package main

import (
	"fmt"
	"sync"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	workers  = kingpin.Flag("workers", "Number of workers").Short('w').Default("5").Int()
	batch    = kingpin.Flag("batch", "Batch size for database commit").Short('b').Default("1000").Int()
	host     = kingpin.Flag("host", "Postgres hostname").Required().String()
	user     = kingpin.Flag("user", "Postgres username").Default("postgres").String()
	password = kingpin.Flag("password", "Postgres password").String()
	db       = kingpin.Flag("db", "Postgres database").Required().String()
)

func main() {
	kingpin.Parse()

	recs := make(chan *Record)

	var wg sync.WaitGroup
	wg.Add(*workers)

	for i := 0; i < *workers; i++ {
		go func() {
			defer wg.Done()

			// separate connection to db for each worker
			api := NewDBApi(*host, *db, *user, *password, *batch)
			defer api.Flush()

			filter := NewFilter(recs, []Processor{
				&ProcessorTemperature{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorFeedingBottle{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorGrowth{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorLactation{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorPumping{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorSleep{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorTemperature{ProcessorBasic: NewProcessorBasic(api)},
				&ProcessorWeight{ProcessorBasic: NewProcessorBasic(api)},
			})
			filter.Run()
		}()
	}

	// connect to db and process each record
	api := NewDBApi(*host, *db, *user, *password, *batch)
	bar := pb.StartNew(api.GetCount())
	api.ForEachRecord(func(rec *Record) {
		bar.Increment()
		recs <- rec
	})

	close(recs)
	wg.Wait()

	fmt.Println("\nDONE!")
}
