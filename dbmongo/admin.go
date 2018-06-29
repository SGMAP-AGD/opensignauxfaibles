package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// AdminID Collection key
type AdminID struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}

// Admin stockage des parametres admin
type Admin struct {
	ID     AdminID     `json:"id" bson:"_id"`
	Params interface{} `json:"params" bson:"params"`
}

func (batch *Admin) load(batchKey string, db *mgo.Database) error {
	err := db.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

func (batch *Admin) save(db *mgo.Database) error {
	_, err := db.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

func newBatch(c *gin.Context) {
	var admin Admin
	db := c.Keys["DB"].(*mgo.Database)

	admin.ID.Key = c.Params.ByName("batchID")
	admin.ID.Type = "batch"

	err := admin.save(db)

	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, admin)
	}
}

func listBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch []Admin
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.batch").All(&batch)
	c.JSON(200, batch)
}

func cloneDB(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)

	from := viper.GetString("DB")
	to := c.Params.ByName("to")

	var result interface{}
	declareDatabaseCopy(db, from, to)
	err := db.Run(bson.M{"eval": "copyDatabase()"}, result)
	if err != nil {
		c.JSON(500, err)
		fmt.Println(err)
		return
	}
	removeDatabaseCopy(db)
	c.JSON(200, err)
}
