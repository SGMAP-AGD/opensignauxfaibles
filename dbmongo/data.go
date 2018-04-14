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
