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

	// mapFct, _ := ioutil.ReadFile("js/dataDebit_map.js")
	// reduceFct, _ := ioutil.ReadFile("js/dataDebit_reduce.js")
	// finalizeFct, _ := ioutil.ReadFile("js/dataDebit_finalize.js")

	// job := &mgo.MapReduce{
	// 	Map:      string(mapFct),
	// 	Reduce:   string(reduceFct),
	// 	Finalize: string(finalizeFct),
	// }

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).One(&data)
	c.JSON(200, data)
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/algo1_map.js")
	reduceFct, _ := ioutil.ReadFile("js/algo1_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/algo1_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	// var etablissement []struct {
	// 	ID    string      `json:"id" bson:"_id"`
	// 	Value interface{} `json:"value" bson:"value"`
	// }
	var etablissement []interface{}

	// db.C("Etablissement").Pipe(
	// 	[]bson.M{
	// 		{"$match": bson.M{"value.siret": c.Params.ByName("siret")}},
	// 		{"$addFields": bson.M{"value.siren": bson.M{"$substr": []interface{}{"$value.siret", 0, 9}}}},
	// 		{"$lookup": bson.M{"from": "Entreprise",
	// 			"localField":   "value.siren",
	// 			"foreignField": "value.siren",
	// 			"as":           "value.entreprise",
	// 		},
	// 		},
	// 	}).Batch

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceAll(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/algo1_map.js")
	reduceFct, _ := ioutil.ReadFile("js/algo1_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/algo1_finalize.js")

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

func compact(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/compact_map.js")
	reduceFct, _ := ioutil.ReadFile("js/compact_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/compact_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var etablissement []interface{}

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func compactAll(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/compact_map.js")
	reduceFct, _ := ioutil.ReadFile("js/compact_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/compact_finalize.js")

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

func compactEntreprise(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/compact_map_entreprise.js")
	reduceFct, _ := ioutil.ReadFile("js/compact_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/compact_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var entreprise []interface{}

	db.C("Entreprise").Find(bson.M{"value.siren": c.Params.ByName("siren")}).MapReduce(job, &entreprise)

	c.JSON(200, entreprise)
}

func compactAllEntreprise(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/compact_map_entreprise.js")
	reduceFct, _ := ioutil.ReadFile("js/compact_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/compact_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "Entreprise"},
	}

	db.C("Entreprise").Find(nil).MapReduce(job, nil)

	c.JSON(200, "OK dude")
}
