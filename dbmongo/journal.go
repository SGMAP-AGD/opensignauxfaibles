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

type journalChannel chan journalEvent

var journalClientChannels = []journalChannel{}
var mainJournalChannel = journalDispatch()
var addClientChannel = make(chan journalChannel)

var debug = journalPriority("debug")
var info = journalPriority("info")
var warning = journalPriority("warning")
var critical = journalPriority("critical")

// journalAddClientChannel surveille l'ajout de nouveaux clients pour propager le channel d'events
func journalAddClient() {
	for clientChannel := range addClientChannel {
		journalClientChannels = append(journalClientChannels, clientChannel)
	}
}

// journal dispatch un event vers les clients et l'enregistre dans la bdd
func journalDispatch() chan journalEvent {
	channel := make(journalChannel)
	go func() {
		for event := range channel {
			for _, clientChannel := range journalClientChannels {
				clientChannel <- event
			}
		}
	}()
	return channel
}

func log(priority journalPriority, code journalCode, comment string) journalEvent {
	e := journalEvent{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Comment:  comment,
		Priority: priority,
		Code:     code,
	}
	mainJournalChannel <- e
	return e
}
