package main

import "time"

type Report struct {
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
}

type ReportsByCreated []Report

func (a ReportsByCreated) Len() int           { return len(a) }
func (a ReportsByCreated) Less(i, j int) bool { return a[j].CreatedAt.After(a[i].CreatedAt) }
func (a ReportsByCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
