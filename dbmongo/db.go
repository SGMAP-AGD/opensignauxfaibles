package main

import (
	"log"
	"time"

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

	dbInsertEntreprise := insertEntreprise(db)
	dbInsertEtablissement := insertEtablissement(db)

	go func() {
		for range time.Tick(2 * time.Second) {
			dbInsertEntreprise <- ValueEntreprise{}
			dbInsertEtablissement <- ValueEtablissement{}
		}
	}()

	return func(c *gin.Context) {
		c.Set("insertEntreprise", dbInsertEntreprise)
		c.Set("insertEtablissement", dbInsertEtablissement)
		c.Set("DB", db)
		c.Next()
	}
}

func insertEntreprise(db *mgo.Database) chan ValueEntreprise {
	source := make(chan ValueEntreprise)

	go func(chan ValueEntreprise) {
		buffer := make([]interface{}, 0)
		i := 0

		for value := range source {
			if value.Value.Siren == "" || len(buffer) >= 1000 {
				if i > 0 {
					db.C("Entreprise").Insert(buffer...)
					buffer = make([]interface{}, 0)
					i = 0
				}
			} else {
				value.ID = bson.NewObjectId()
				buffer = append(buffer, value)
				i++
			}

		}
	}(source)

	return source
}

func insertEtablissement(db *mgo.Database) chan ValueEtablissement {
	source := make(chan ValueEtablissement)

	go func(chan ValueEtablissement) {
		buffer := make([]interface{}, 0)
		i := 0

		for value := range source {
			if value.Value.Siret == "" || len(buffer) >= 1000 {
				if i > 0 {
					db.C("Etablissement").Insert(buffer...)
					buffer = make([]interface{}, 0)
					i = 0
				}
			} else {
				value.ID = bson.NewObjectId()
				buffer = append(buffer, value)
				i++
			}

		}
	}(source)

	return source
}
