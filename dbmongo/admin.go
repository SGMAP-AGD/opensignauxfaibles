package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

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
	ID       AdminID    `json:"id" bson:"_id"`
	Files    BatchFiles `json:"files" bson:"files"`
	Readonly bool       `json:"readonly" bson:"readonly"`
	Params   struct {
		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
	} `json:"params" bson:"param"`
	ChanEtablissement chan *ValueEtablissement `json:"-" bson:"-"`
	ChanEntreprise    chan *ValueEntreprise    `json:"-" bson:"-"`
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

func updateBatch(c *gin.Context) {
	batch := AdminBatch{}
	err := c.Bind(&batch)
	spew.Dump(err)
	if err != nil {
		c.JSON(500, err)
		return
	}

	db := c.Keys["DB"].(*mgo.Database)

	err = batch.save(db)
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement")
		return
	}
	c.JSON(200, batch)
}

func listBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch []AdminBatch
	err := db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
	if err != nil {
		spew.Dump(err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, batch)
}

func getBatchesID(db *mgo.Database) []string {
	var batch []AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
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
	c.JSON(200, []struct {
		Type    string `json:"type" bson:"type"`
		Libelle string `json:"text" bson:"text"`
	}{
		{"admin_urssaf", "Siret/Compte URSSAF"},
		{"apconso", "Consommation Activité Partielle"},
		{"bdf", "Ratios Banque de France"},
		{"cotisation", "Cotisations URSSAF"},
		{"delai", "Délais URSSAF"},
		{"dpae", "Déclaration Préalable à l'embauche"},
		{"interim", "Base Interim"},
		{"altares", "Base Altarès"},
		{"apdemande", "Demande Activité Partielle"},
		{"ccsf", "Stock CCSF à date"},
		{"debit", "Débits URSSAF"},
		{"dmmo", "Déclaration Mouvement de Main d'Œuvre"},
		{"effectif", "Emplois URSSAF"},
		{"sirene", "Base GéoSirene"},
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

// NAF libellés et liens N5/N1
type NAF struct {
	N1    map[string]string `json:"n1" bson:"n1"`
	N5    map[string]string `json:"n5" bson:"n5"`
	N5to1 map[string]string `json:"n5to1" bson:"n5to1"`
}

func loadNAF() (NAF, error) {
	naf := NAF{}
	naf.N1 = make(map[string]string)
	naf.N5 = make(map[string]string)
	naf.N5to1 = make(map[string]string)

	NAF1 := viper.GetString("NAF_L1")
	NAF5 := viper.GetString("NAF_L5")
	NAF5to1 := viper.GetString("NAF_5TO1")

	NAF1File, NAF1err := os.Open(NAF1)
	if NAF1err != nil {
		return NAF{}, NAF1err
	}

	NAF1reader := csv.NewReader(bufio.NewReader(NAF1File))
	NAF1reader.Comma = ';'
	NAF1reader.Read()
	for {
		row, error := NAF1reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		naf.N1[row[0]] = row[1]
		fmt.Println(row)
	}

	NAF5to1File, NAF5to1err := os.Open(NAF5to1)
	if NAF5to1err != nil {
		return NAF{}, NAF1err
	}

	NAF5to1reader := csv.NewReader(bufio.NewReader(NAF5to1File))
	NAF5to1reader.Comma = ';'
	NAF5to1reader.Read()
	for {
		row, error := NAF5to1reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		naf.N5to1[row[0]] = row[1]
	}

	NAF5File, NAF5err := os.Open(NAF5)
	if NAF5err != nil {
		return NAF{}, NAF1err
	}

	NAF5reader := csv.NewReader(bufio.NewReader(NAF5File))
	NAF5reader.Comma = ';'
	NAF5reader.Read()
	for {
		row, error := NAF5reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		naf.N5[row[0]] = row[1]

	}

	return naf, nil
}
