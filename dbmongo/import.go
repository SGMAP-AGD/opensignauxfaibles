package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

func adminFiles(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := listFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

type fileSummary struct {
	Name string    `json:"name" bson:"name"`
	Size int64     `json:"size" bson:"size"`
	Date time.Time `json:"date" bson:"date"`
}

func listFiles(basePath string) ([]fileSummary, error) {
	var files []fileSummary
	basePathConf := viper.GetString("APP_DATA")
	b := len(basePathConf)

	currentFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []fileSummary{}, err
	}

	for _, file := range currentFiles {
		if file.IsDir() {
			subPath := fmt.Sprintf("%s/%s", basePath, file.Name())
			subFiles, err := listFiles(subPath)
			if err != nil {
				return []fileSummary{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, fileSummary{
				Name: fmt.Sprintf("%s/%s", basePath, file.Name())[b:],
				Size: file.Size(),
				Date: file.ModTime(),
			})
		}
	}

	return files, nil
}

var importFunctions = map[string]func(*AdminBatch) error{
	"apconso":    importAPConso,
	"bdf":        importBDF,
	"diane":      importDiane,
	"cotisation": importCotisation,
	"delai":      importDelai,
	"dpae":       importDPAE,
	"altares":    importAltares,
	"apdemande":  importAPDemande,
	"ccsf":       importCCSF,
	"debit":      importDebit,
	"effectif":   importEffectif,
	"sirene":     importSirene,
}

func purge(c *gin.Context) {
	db, _ := c.Keys["db"].(*mgo.Database)
	db.C("Etablissement").RemoveAll(nil)
	db.C("Entreprise").RemoveAll(nil)
	c.String(200, "Done")
}

func importBatch(c *gin.Context) {
	batch := AdminBatch{}
	batchKey := c.Params.ByName("batch")
	db := c.Keys["db"].(*mgo.Database)
	chanEtablissement := c.Keys["ChanEtablissement"].(chan *ValueEtablissement)
	chanEntreprise := c.Keys["ChanEntreprise"].(chan *ValueEntreprise)
	batch.load(batchKey, db, chanEtablissement, chanEntreprise)

	if !batch.Readonly {
		for _, fn := range importFunctions {
			err := fn(&batch)
			if err != nil {
				c.JSON(500, err)
				return
			}
		}
	} else {
		c.JSON(403, "Ce lot est fermé, import impossible.")
	}
}

func lastOpenBatch(c *gin.Context) {
	db := c.Keys["db"].(*mgo.Database)
	batches := getBatches(db)
	spew.Dump(batches)
}
