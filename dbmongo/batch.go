package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

// AdminID Collection key
type AdminID struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}

// AdminBatch décrit les informations d'un batch
type AdminBatch struct {
	ID          AdminID             `json:"id" bson:"_id"`
	Description string              `json:"description" bson:"description"`
	Files       map[string][]string `json:"files" bson:"files"`
}

func getNewBatch(c *gin.Context) {

}
func newBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch AdminBatch
	err := c.Bind(batch)
	if err != nil {
		c.JSON(500, err)
		return
	}

	info, err := db.C("Admin").Upsert(batch.ID, batch)

	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, info)
	}
}
