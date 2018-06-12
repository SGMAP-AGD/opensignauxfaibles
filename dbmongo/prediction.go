package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func predictionBrowse(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	batch := c.Params.ByName("batch")
	algo := c.Params.ByName("algo")
	paramPage := c.Params.ByName("page")

	page, err := strconv.Atoi(paramPage)

	if err != nil {
		c.JSON(500, err)
		return
	}

	//siret := c.Params.ByName("siret")
	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$match": bson.M{"_id.batch": batch, "_id.algo": algo}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"score": -1}})
	pipeline = append(pipeline, bson.M{"$skip": 50 * page})
	pipeline = append(pipeline, bson.M{"$limit": 50})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
	}})
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Etablissement",
		"localField":   "_id.siret",
		"foreignField": "_id",
		"as":           "etablissement"}})
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Entreprise",
		"localField":   "siren",
		"foreignField": "_id",
		"as":           "entreprise"}})
	var result []interface{}
	err = db.C("Prediction").Pipe(pipeline).All(&result)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, result)
}
