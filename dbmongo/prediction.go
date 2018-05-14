package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

// Prediction type pour stocker les prédictions période par période
type Prediction struct {
	RandomForest float64 `json:"random_forest" bson:"random_forest"`
}

func injectPrediction(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	prediction := []struct {
		ID          string    `json:"id" bson:"_id"`
		Periode     time.Time `json:"periode" bson:"periode"`
		Probabilite float64   `json:"prob" bson:"prob"`
		Siret       string    `json:"Siret"`
	}{}

	db.C("prediction").Find(nil).All(prediction)

	c.JSON(200, nil)
}
