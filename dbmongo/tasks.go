package main

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func getTasks(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	id := claims[identityKey]

	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"scope": id,
	}})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"date": 1,
	}})

	pipeline = append(pipeline, bson.M{"$group": bson.M{
		"_id":       "$_id.siret",
		"lastDate":  bson.M{"$max": "$_id.date"},
		"firstDate": bson.M{"$min": "$_id.date"},
		"tasks":     bson.M{"$push": "$$ROOT"},
	}})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"firstDate": -1,
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "PublicEtablissement",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "etablissement",
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "PublicEntreprise",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "entreprise",
	}})

	var result []interface{}
	err := db.DB.C("Tasks").Pipe(pipeline).All(&result)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, result)
}
