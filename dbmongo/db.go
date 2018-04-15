package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DB Initialisation de la connexion MongoDB
func DB() gin.HandlerFunc {

	mongodb, err := mgo.Dial("127.0.0.1")
	db := mongodb.DB("jason")

	// pousse des fonctions partagées JS
	declareServerFunctions(db)
	if err != nil {
		log.Panic(err)
	}

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func insertValue(db *mgo.Database, value Value) {
	value.ID = bson.NewObjectId()
	db.C("Etablissement").Insert(value)
}
