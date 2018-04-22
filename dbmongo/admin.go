package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Path Chemins d'accès aux répertoires sources
type Path struct {
	Debit      string `json:"debit" bson:"debit"`
	Delais     string `json:"delais" bson:"delais"`
	Cotisation string `json:"cotisation" bson:"cotisation"`
	Ccsf       string `json:"ccsf" bson:"ccsf"`
	Altares    string `json:"altares" bson:"altares"`
	APDemande  string `json:"apdemande" bson:"apdemande"`
	APConso    string `json:"apconso" bson:"apconso"`
	BDF        string `json:"bdf" bson:"bdf"`
}

// Region Descripteur de région
type Region struct {
	ID    string `json:"id" bson:"_id"`
	Label string `json:"label" bson:"label"`
}

// AdminRegion Lister les régions
func AdminRegion(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	var regions []Region
	db.C("region").Find(nil).All(&regions)

	c.JSON(200, regions)
}

// AdminRegionAdd Ajouter une région
func AdminRegionAdd(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var region Region
	err := c.BindJSON(&region)

	if err != nil {
		c.JSON(500, "not that good")
		log.Panic(err)
	} else {
		chg, err := db.C("region").Upsert(
			bson.M{"id": region.ID},
			region,
		)
		if err != nil {
			fmt.Println(err, chg)
		}
		c.JSON(200, "ok")
	}
}

// AdminRegionDelete Supprimer une région
func AdminRegionDelete(c *gin.Context) {
	c.JSON(200, "ok")
}
