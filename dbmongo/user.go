package main

import (
	"github.com/globalsign/mgo/bson"
)

// User identification d'utilisateur
type User struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Permission []string      `json:"permission" bson:"permission"`
	Nom        string        `json:"nom" bson:"nom"`
	Prenom     string        `json:"prenom" bson:"prenom"`
	Contact    string        `json:"telephone" bson:"telephone"`
	Login      string        `json:"login" bson:"login"`
	Pass       string        `json:"password" bson:"password"`
	Region     bson.ObjectId `json:"region_id" bson:"region_id"`
}

// Region Détermine les propriétés d'une région
type Region struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Path  string        `json:"path" bson:"path"`
	Label string        `json:"label" bson:"label"`
}
