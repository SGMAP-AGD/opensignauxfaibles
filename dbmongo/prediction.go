package main

import (
	"fmt"
	"time"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

// Prediction type pour stocker les prédictions période par période
type Prediction struct {
	RandomForest float64 `json:"random_forest" bson:"random_forest"`
}

func injectPrediction(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")

	predictions := []struct {
		ID          string    `json:"id" bson:"_id"`
		Periode     time.Time `json:"periode" bson:"periode"`
		Probabilite float64   `json:"prob" bson:"prob"`
		Siret       string    `json:"Siret"`
	}{}

	db.C("prediction").Find(nil).All(&predictions)

	for _, p := range predictions {
		prediction := Prediction{
			RandomForest: p.Probabilite,
		}

		hash := fmt.Sprintf("%x", structhash.Md5(prediction, 1))

		etablissement := Value{
			Value: Etablissement{
				Siret:  p.Siret,
				Region: region,
				Batch: map[string]Batch{
					batch: Batch{
						Compact: map[string]bool{
							"status": false,
						},
						Prediction: map[string]Prediction{
							hash: prediction,
						},
					},
				},
			},
		}

		db.C("Etablissement").Insert(&etablissement)
	}
	c.JSON(200, nil)
}
