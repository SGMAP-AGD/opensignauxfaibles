package main

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func dataDebit(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var data interface{}

	mapFct, _ := ioutil.ReadFile("js/dataDebit_map.js")
	reduceFct, _ := ioutil.ReadFile("js/dataDebit_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/dataDebit_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).MapReduce(job, &data)

	c.JSON(200, data)
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map2.js")
	reduceFct, _ := ioutil.ReadFile("reduce2.js")
	finalizeFct, _ := ioutil.ReadFile("finalize2.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var etablissement []struct {
		ID    string      `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(bson.M{"_id": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceAll(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map2.js")
	reduceFct, _ := ioutil.ReadFile("reduce2.js")
	finalizeFct, _ := ioutil.ReadFile("finalize2.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "algo1"},
	}

	var etablissement []struct {
		ID    string      `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(bson.M{"value.index.algo1": true}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func debug(c *gin.Context) {
	ursaff := c.Params.ByName("urssaf")
	date, _ := UrssafToPeriod(ursaff)
	c.JSON(200, date)
}

func browseEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var etablissement []interface{}
	db.C("testcollection").Find(bson.M{"_id": c.Params.ByName("siret")}).All(&etablissement)
	c.JSON(200, etablissement)
}

func browseOrig(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var etablissement []interface{}
	db.C("Etablissement").Find(bson.M{"siret": c.Params.ByName("siret")}).All(&etablissement)
	c.JSON(200, etablissement)
}

func reduceEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map.js")
	reduceFct, _ := ioutil.ReadFile("reduce.js")
	finalizeFct, _ := ioutil.ReadFile("finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var etablissement []interface{}

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceEtablissements(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map.js")
	reduceFct, _ := ioutil.ReadFile("reduce.js")
	finalizeFct, _ := ioutil.ReadFile("finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "Etablissement"},
	}

	var etablissement []struct {
		ID    string        `json:"id" bson:"_id"`
		Value Etablissement `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(nil).MapReduce(job, nil)

	c.JSON(200, etablissement)
}
