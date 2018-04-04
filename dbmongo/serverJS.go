package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ServerJSFunc Function à injecter dans l'instance MongoDB
type ServerJSFunc struct {
	ID    string `json:"id" bson:"_id"`
	Value string `json:"value" bson:"value"`
}

// Add Méthode pour insérer
func (f ServerJSFunc) Add(db *mgo.Database) {
	db.C("system.js").Upsert(bson.M{"_id": f.ID}, f)
}
