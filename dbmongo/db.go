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
	mongodb.SetSocketTimeout(3600 * time.Second)
	db := mongodb.DB(dbDatabase)

	// pousse les fonctions partagées JS
	declareServerFunctions(db)

	if err != nil {
		log.Panic(err)
	}

	dbInsertEntreprise := insertEntreprise(db)
	dbInsertEtablissement := insertEtablissement(db)

	go func() {
		for range time.Tick(30 * time.Second) {
			dbInsertEntreprise <- &ValueEntreprise{}
			dbInsertEtablissement <- &ValueEtablissement{}
		}
	}()

	return func(c *gin.Context) {
		c.Set("insertEntreprise", dbInsertEntreprise)
		c.Set("insertEtablissement", dbInsertEtablissement)
		c.Set("DB", db)
		c.Next()
	}
}

func insertEntreprise(db *mgo.Database) chan *ValueEntreprise {
	source := make(chan *ValueEntreprise, 1000)

	go func(chan *ValueEntreprise) {
		buffer := make(map[string]*ValueEntreprise)
		objects := make([]interface{}, 0)
		i := 0

		for value := range source {
			if value.Value.Siren == "" || i >= 100 {
				for _, v := range buffer {
					objects = append(objects, *v)
				}
				db.C("Entreprise").Insert(objects...)

				buffer = make(map[string]*ValueEntreprise)
				objects = make([]interface{}, 0)
				i = 0
			} else {
				if knowValue, ok := buffer[value.Value.Siren]; ok {
					newValue, _ := (*knowValue).merge(*value)
					buffer[value.Value.Siren] = &newValue
				} else {
					value.ID = bson.NewObjectId()
					buffer[value.Value.Siren] = value
					i++
				}
			}

		}
	}(source)

	return source
}

func insertEtablissement(db *mgo.Database) chan *ValueEtablissement {
	source := make(chan *ValueEtablissement, 1000)

	go func(chan *ValueEtablissement) {
		buffer := make(map[string]*ValueEtablissement)
		objects := make([]interface{}, 0)
		i := 0

		for value := range source {
			if value.Value.Siret == "" || i >= 100 {
				for _, v := range buffer {
					objects = append(objects, *v)
				}
				go func(o []interface{}) { db.C("Etablissement").Insert(o...) }(objects)

				buffer = make(map[string]*ValueEtablissement)
				objects = make([]interface{}, 0)
				i = 0
			} else {
				if knowValue, ok := buffer[value.Value.Siret]; ok {
					newValue, _ := (*knowValue).merge(*value)
					buffer[value.Value.Siret] = &newValue
				} else {
					value.ID = bson.NewObjectId()
					buffer[value.Value.Siret] = value
					i++
				}
			}

		}
	}(source)

	return source
}
