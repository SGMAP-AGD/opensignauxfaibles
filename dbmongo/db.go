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

	dbWorker := insertWorker(db)

	go func() {
		for range time.Tick(2 * time.Second) {
			dbWorker <- Value{}
		}
	}()

	return func(c *gin.Context) {
		c.Set("DBW", dbWorker)
		c.Set("DB", db)
		c.Next()
	}
}

func insertWorker(db *mgo.Database) chan Value {
	source := make(chan Value)

	go func(chan Value) {
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

// func insertValue(db *mgo.Database, value Value) error {
// 	if value.Value.Siret != "" {
// 		value.ID = bson.NewObjectId()
// 		err := db.C("Etablissement").Insert(value)
// 		return err
// 	}
// 	return nil
// }

// func insertValueEntreprise(db *mgo.Database, value ValueEntreprise) error {
// 	if value.Value.Siren != "" {
// 		value.ID = bson.NewObjectId()
// 		err := db.C("Entreprise").Insert(value)
// 		return err
// 	}
// 	return nil
// }
