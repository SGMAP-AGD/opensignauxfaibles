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
		"delai",
		"effectif",
		"sirene",
		"bdf",
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
func GetFileList(basePath string, period string) (map[string][]string, map[string]error) {
	list := make(map[string][]string)
	var l []os.FileInfo
	err := make(map[string]error)
	directories := []string{
		"admin_urssaf",
		"altares",
		"apdemande",
		"apconso",
		"bdf",
		"ccsf",
		"cotisation",
		"debit",
		"delai",
		"effectif",
		"sirene",
		"bdf",
	}

	for _, dir := range directories {
		l, err[dir] = ioutil.ReadDir(fmt.Sprintf("%s/%s/%s", basePath, period, dir))
		for _, f := range l {
			list[dir] = append(list[dir], fmt.Sprintf("%s/%s/%s/%s", basePath, period, dir, f.Name()))
		}
	}
	return list, err
}

func importAll(c *gin.Context) {
	importAltares(c)
	importAPConso(c)
	importAPDemande(c)
	importEffectif(c)
	importDebit(c)
	importBDF(c)
	importCotisation(c)
	importDelai(c)
	importSirene(c)
}

func purge(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("Etablissement").RemoveAll(nil)
	db.C("Entreprise").RemoveAll(nil)
	c.String(200, "Done")
}
