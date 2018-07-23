package main

import (
	"errors"
	"fmt"
	"time"

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
	ID                AdminID                      `json:"id" bson:"_id"`
	Files             BatchFiles                   `json:"files" bson:"files"`
	Open              bool                         `json:"open" bson:"open"`
	Draft             bool                         `json:"draft" bson:"draft"`
	DateDebut         time.Time                    `json:"date_debut" bson:"date_debut"`
	DateFin           time.Time                    `json:"date_fin" bson:"date_fin"`
	DateFinEffectif   time.Time                    `json:"date_fin_effectif" bson:"date_fin_effectif"`
	Params            map[string]map[string]string `json:"params" bson:"params"`
	ChanEtablissement chan *ValueEtablissement     `json:"-" bson:"-"`
	ChanEntreprise    chan *ValueEntreprise        `json:"-" bson:"-"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

func (batchFiles BatchFiles) attachFile(fileType string, file string) {
	batchFiles[fileType] = append(batchFiles[fileType], file)
}

func (batch *AdminBatch) load(
	batchKey string,
	db *mgo.Database,
	chanEtablissement chan *ValueEtablissement,
	chanEntreprise chan *ValueEntreprise) error {
	err := db.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	batch.ChanEntreprise = chanEntreprise
	batch.ChanEtablissement = chanEtablissement
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
	chanEtablissement := c.Keys["ChanEtablissement"].(chan *ValueEtablissement)
	chanEntreprise := c.Keys["ChanEntreprise"].(chan *ValueEntreprise)

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
	err = batch.load(params.Batch, db, chanEtablissement, chanEntreprise)

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

func updateBatch(c *gin.Context) {
	batch := AdminBatch{}
	err := c.Bind(batch)
	if err != nil {
		c.JSON(500, "Valeur de batch non autorisée")
		return
	}

	db := c.Keys["DB"].(*mgo.Database)

	batch.save(db)
}
func listBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch []AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batch)
	c.JSON(200, batch)
}

func getBatchesID(db *mgo.Database) []string {
	var batch []AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batch)
	var batchesID []string
	for _, b := range batch {
		batchesID = append(batchesID, b.ID.Key)
	}
	return batchesID
}

func getBatches(db *mgo.Database) []*AdminBatch {
	var batches []*AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches
}

func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
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
