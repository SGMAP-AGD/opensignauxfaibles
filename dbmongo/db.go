package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DB Initialisation de la connexion MongoDB
func DB() gin.HandlerFunc {

	dbDial := viper.GetString("DB_DIAL")
	dbDatabase := viper.GetString("DB")

	mongodb, err := mgo.Dial(dbDial)
	db := mongodb.DB(dbDatabase)

	// pousse les fonctions partagées JS
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
	if value.Value.Siret != "" {
		value.ID = bson.NewObjectId()
		db.C("Etablissement").Insert(value)
	}
}

func insertValueEntreprise(db *mgo.Database, value ValueEntreprise) {
	if value.Value.Siren != "" {
		value.ID = bson.NewObjectId()
		db.C("Entreprise").Insert(value)
	}
}
