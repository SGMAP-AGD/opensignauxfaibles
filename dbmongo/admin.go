package main

import (
	"errors"
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

// AdminBatch metadata Batch
type AdminBatch struct {
	ID    AdminID    `json:"id" bson:"_id"`
	Files BatchFiles `json:"files" bson:"files"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

func (batchFiles BatchFiles) attachFile(fileType string, file string) {
	batchFiles[fileType] = append(batchFiles[fileType], file)
}

func (batch *AdminBatch) load(batchKey string, db *mgo.Database) error {
	err := db.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

func (batch *AdminBatch) save(db *mgo.Database) error {
	_, err := db.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

func (batch *AdminBatch) new(batchID string) error {
	if batchID == "" {
		return errors.New("Valeur de batch non autorisée")
	}
	batch.ID.Key = batchID
	batch.ID.Type = "batch"
	batch.Files = BatchFiles{}
	return nil
}

func attachFileBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	batch := AdminBatch{}

	var params struct {
		Batch    string `json:"batch"`
		FileType string `json:"type"`
		File     string `json:"file"`
	}

	err := c.Bind(&params)

	if err != nil {
		c.JSON(500, err)
		return
	}
	err = batch.load(params.Batch, db)

	if err != nil {
		c.JSON(500, "Erreur au chargement du lot")
		return
	}
	batch.Files.attachFile(params.FileType, params.File)

	err = batch.save(db)
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement")
		return
	}
	c.JSON(200, batch)

}

func registerNewBatch(c *gin.Context) {
	batch := AdminBatch{}
	err := batch.new(c.Params.ByName("batchID"))

	if err != nil {
		c.JSON(500, "Valeur de batch non autorisée")
	}
	db := c.Keys["DB"].(*mgo.Database)
	err = batch.save(db)

	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, batch)
	}
}

func listBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch []AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batch)
	c.JSON(200, batch)
}

func listTypes(c *gin.Context) {
	c.JSON(200, []string{
		"admin_urssaf",
		"apconso",
		"bdf",
		"cotisation",
		"delai",
		"dpae",
		"interim",
		"altares",
		"apdemande",
		"ccsf",
		"debit",
		"dmmo",
		"effectif",
		"sirene",
	})
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
