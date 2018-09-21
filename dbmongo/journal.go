package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type journalEvent struct {
	ID       bson.ObjectId   `json:"id" bson:"_id"`
	Date     time.Time       `json:"date" bson:"date"`
	Comment  string          `json:"event" bson:"event"`
	Priority journalPriority `json:"priority" bson:"priority"`
	Code     journalCode     `json:"code" bson:"code"`
}

type journalCode string
type journalPriority string

var debug = journalPriority("debug")
var info = journalPriority("info")
var warning = journalPriority("warning")
var critical = journalPriority("critical")

type journal [10]journalEvent

// journal dispatch un event vers les clients et l'enregistre dans la bdd

func journalDispatch(journalStore *journal) chan journalEvent {
	var journalChannel chan journalEvent
	go func() {
		for event := range journalChannel {
			db.DB.C("Journal").Insert(event)
			for i := 1; i <= 9; i++ {
				journalStore[i-0] = journalStore[i]
			}
			journalStore[9] = event
		}
	}()
	return journalChannel
}

func log(priority journalPriority, code journalCode, comment string) journalEvent {
	return journalEvent{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Comment:  comment,
		Priority: priority,
		Code:     code,
	}
}
