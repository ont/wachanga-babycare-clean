package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBApi struct {
	db *sqlx.DB

	ch chan *Record

	lock       sync.Mutex
	bsize, cnt int
	tx         *sqlx.Tx
}

func NewDBApi(host, dbname, user, password string, bsize int) *DBApi {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable",
		user, password, host, dbname,
	))

	if err != nil {
		log.Fatal(err)
	}

	return &DBApi{
		db:    db,
		ch:    make(chan *Record),
		bsize: bsize,
		tx:    db.MustBegin(),
	}
}

func (d *DBApi) GetCount() (cnt int) {
	err := d.db.Get(&cnt, "SELECT count(*) FROM event")
	if err != nil {
		log.Fatalln("Can't count events:", err)
	}
	return
}

func (d *DBApi) ForEachRecord(callback func(*Record)) {

	rows, err := d.db.Queryx("SELECT id, event_type, raw_value FROM event")
	if err != nil {
		log.Fatalln("Can't select events:", err)
	}

	rec := Record{}
	for rows.Next() {
		err := rows.StructScan(&rec)
		if err != nil {
			log.Fatalln(err)
		}

		recCopy := rec
		recCopy.Parse() // prepare JsonValue field
		callback(&recCopy)
	}
}

func (d *DBApi) Delete(rec *Record) {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.batchFlush()

	//fmt.Println("Deleting", rec.Id)
	_, err := d.tx.NamedExec("DELETE FROM event WHERE id=:id", rec)
	if err != nil {
		log.Fatalln(err)
	}
}

func (d *DBApi) Save(rec *Record) {
	d.lock.Lock()
	defer d.lock.Unlock()
	defer d.batchFlush()

	rec.Serialize()

	//fmt.Println("Saving", rec.Id)
	_, err := d.tx.NamedExec("UPDATE event SET raw_value=:raw_value WHERE id=:id", rec)
	if err != nil {
		log.Fatalln(err)
	}
}

func (d *DBApi) batchFlush() {
	d.cnt++

	if d.cnt > d.bsize {
		d.Flush()
		d.cnt = 0
	}
}

// flush immediatelly
func (d *DBApi) Flush() {
	err := d.tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	d.tx = d.db.MustBegin()
}
