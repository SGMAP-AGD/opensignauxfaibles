package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
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
	//"apconso":    importAPConso,
	//"bdf":        importBDF,
	// "diane":      importDiane,
	// "cotisation": importCotisation,
	"delai": importDelai,
	// "dpae":       importDPAE,
	// "altares":    importAltares,
	// "apdemande":  importAPDemande,
	// "ccsf":       importCCSF,
	// "debit":      importDebit,
	// "effectif":   importEffectif,
	// "sirene":     importSirene,
}

func purge(c *gin.Context) {
	db.DB.C("Etablissement").RemoveAll(nil)
	db.DB.C("Entreprise").RemoveAll(nil)
	c.String(200, "Done")
}

func importBatchHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batch")
	batch := AdminBatch{}
	batch.load(batchKey)
	go importBatch(&batch)
}

func importBatch(batch *AdminBatch) {
	if !batch.Readonly {
		for _, fn := range importFunctions {
			fmt.Println(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
			err := fn(batch)
			if err != nil {
				journal("Erreur à l'import du fichier: "+err.Error(), "Critique")
			}
		}
	} else {
		journal("Ce lot est fermé, import impossible.", "Critique")
	}
}

func journal(event string, priority string) {
	entry := struct {
		ID       bson.ObjectId `json:"id" bson:"_id"`
		Date     time.Time     `json:"date" bson:"date"`
		Event    string        `json:"event" bson:"event"`
		Priority string        `json:"priority" bson:"priority"`
	}{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Event:    event,
		Priority: priority,
	}
	db.DB.C("Journal").Insert(entry)
}
