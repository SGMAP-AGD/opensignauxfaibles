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
	pipeline = append(pipeline, bson.M{"$skip": 10 * page})
	pipeline = append(pipeline, bson.M{"$limit": 10})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Etablissement",
		"localField":   "_id.siret",
		"foreignField": "_id",
		"as":           "etablissement"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"etablissement": bson.M{"$arrayElemAt": []interface{}{"$etablissement", 0}}}})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"etablissement": "$etablissement.value"}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Entreprise",
		"localField":   "siren",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"entreprise": bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}}}})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"entreprise": "$entreprise.value"}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from": "Features",
		"let": bson.M{"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
			"siret": "$_id.siret"},
		"pipeline": []interface{}{
			bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
				bson.M{"$eq": []interface{}{"$_id.siren", "$$siren"}},
				bson.M{"$eq": []interface{}{"$_id.batch", batch}},
				bson.M{"$eq": []interface{}{"$_id.algo", algo}},
			}}}},
			bson.M{"$addFields": bson.M{
				"features": bson.M{"$filter": bson.M{"input": "$value",
					"cond": bson.M{
						"$eq": []interface{}{"$$this.siret", "$$siret"}}}}}},
			bson.M{"$project": bson.M{
				"features": "$features",
			}},
		},
		"as": "features"}})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"features": bson.M{"$arrayElemAt": []interface{}{"$features", 0}}}})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"features": "$features.features"}})

	var result []interface{}
	err = db.C("Prediction").Pipe(pipeline).All(&result)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, result)
}
