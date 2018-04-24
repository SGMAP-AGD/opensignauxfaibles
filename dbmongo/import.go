package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

func createRepo(c *gin.Context) {
	// db := c.Keys["DB"].(*mgo.Database)
	basePath := viper.GetString("APP_DATA")

	// db.C("region").Find(bson.M{"ID"})
	directories := []string{
		"admin_urssaf",
		"apconso",
		"apdemande",
		"bdf",
		"ccsf",
		"cotisation",
		"debit",
		"delais",
		"effectif",
	}

	var response map[string]string
	var status int
	for _, directory := range directories {
		path := basePath + "/" + directory
		err := os.MkdirAll(path, 700)
		status = 200
		if err != nil {
			status = 207
			response[path] = err.Error()
		} else {
			response[path] = "ok"
		}
	}
	c.JSON(status, "ok")
}

// GetFileList construit la liste des fichiers à traiter
func GetFileList(basePath string, region string, period string) (map[string][]string, map[string]error) {
	list := make(map[string][]string)
	var l []os.FileInfo
	err := make(map[string]error)
	directories := []string{"admin_urssaf", "altares", "altares", "apdemande", "apconso", "bdf", "ccsf", "cotisation", "debit", "delais", "effectif"}

	for _, dir := range directories {
		l, err[dir] = ioutil.ReadDir(fmt.Sprintf("%s/%s/%s/%s", basePath, region, period, dir))
		for _, f := range l {
			list[dir] = append(list[dir], fmt.Sprintf("%s/%s/%s/%s/%s", basePath, region, period, dir, f.Name()))
		}
	}

	return list, err
}

func importAltares(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	altares := files["altares"][0]

	for etablissement := range parseAltares(altares, batch) {
		etablissement.Region = region
		insertValue(db, Value{Value: etablissement})
	}
}

func importEffectif(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	effectif := files["effectif"][0]

	for etablissement := range parseEffectif(effectif, batch) {
		etablissement.Region = region
		insertValue(db, Value{Value: etablissement})
	}
}

func importAPDemande(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	apdemande := files["apdemande"][0]

	for etablissement := range parseAPDemande(apdemande, batch) {
		etablissement.Region = region
		insertValue(db, Value{Value: etablissement})
	}
}

func importAPConso(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	apconso := files["apconso"][0]

	for etablissement := range parseAPConsommation(apconso, batch) {
		etablissement.Region = region
		insertValue(db, Value{Value: etablissement})
	}
}

func importDebit(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	debit := files["debit"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	value := Value{Value: Etablissement{}}
	value.Value.Batch = make(map[string]Batch)
	value.Value.Batch[batch] = Batch{Debit: map[string]Debit{}}

	for etablissement := range parseDebit(debit, batch) {
		etablissement.Region = region
		etablissement.Siret = mapping[etablissement.Key]
		if value.Value.Siret == etablissement.Siret {
			for k, v := range etablissement.Batch[batch].Debit {
				value.Value.Batch[batch].Debit[k] = v
			}
		} else {
			insertValue(db, value)
			value = Value{Value: etablissement}
		}
	}
	insertValue(db, value)
}

func importCotisation(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	cotisation := files["cotisation"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	value := Value{Value: Etablissement{}}
	value.Value.Batch = make(map[string]Batch)
	value.Value.Batch[batch] = Batch{Cotisation: map[string]Cotisation{}}

	for etablissement := range parseCotisation(cotisation, batch) {
		etablissement.Region = region
		etablissement.Siret = mapping[etablissement.Key]
		if value.Value.Siret == etablissement.Siret {
			for k, v := range etablissement.Batch[batch].Cotisation {
				value.Value.Batch[batch].Cotisation[k] = v
			}
		} else {
			insertValue(db, value)
			value = Value{Value: etablissement}
		}
	}
	insertValue(db, value)
}

func importDelais(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batch := c.Params.ByName("batch")
	region := c.Params.ByName("region")
	files, _ := GetFileList(viper.GetString("APP_DATA"), region, batch)
	delais := files["delais"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	value := Value{Value: Etablissement{}}
	value.Value.Batch = make(map[string]Batch)
	value.Value.Batch[batch] = Batch{Delais: map[string]Delais{}}

	for etablissement := range parseCotisation(delais, batch) {
		etablissement.Region = region
		etablissement.Siret = mapping[etablissement.Key]
		if value.Value.Siret == etablissement.Siret {
			for k, v := range etablissement.Batch[batch].Delais {
				value.Value.Batch[batch].Delais[k] = v
			}
		} else {
			insertValue(db, value)
			value = Value{Value: etablissement}
		}
	}
	insertValue(db, value)
}

func purge(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("Etablissement").RemoveAll(nil)
	c.String(200, "Done")
}
