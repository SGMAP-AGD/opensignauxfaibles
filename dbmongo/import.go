package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

func createRepo(c *gin.Context) {
	// db := c.Keys["DB"].(*mgo.Database)
	basePath := viper.GetString("APP_DATA")
	batch := c.Params.ByName("batch")

	directories := []string{
		"admin_urssaf",
		"altares",
		"apconso",
		"apdemande",
		"bdf",
		"ccsf",
		"cotisation",
		"debit",
		"delai",
		"effectif",
		"sirene",
		"interim",
		"dmmo",
		"dpae",
	}

	response := make(map[string]string)
	var status int
	for _, directory := range directories {
		path := basePath + "/" + batch + "/" + directory
		err := os.MkdirAll(path, 0755)
		status = 200
		if err != nil {
			status = 207
			response[path] = err.Error()
		} else {
			response[path] = "ok"
		}
	}
	c.JSON(status, response)
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
		"interim",
		"dmmo",
		"dpae",
	}

	for _, dir := range directories {
		path := fmt.Sprintf("%s/%s/%s", basePath, period, dir)
		l, err[dir] = ioutil.ReadDir(path)
		for _, f := range l {
			if match, _ := regexp.MatchString("\\.(csv|xls|xlsx)$", f.Name()); match {
				list[dir] = append(list[dir], fmt.Sprintf("%s/%s", path, f.Name()))
			}
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
	importDPAE(c)
}

func purge(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("Etablissement").RemoveAll(nil)
	db.C("Entreprise").RemoveAll(nil)
	c.String(200, "Done")
}
