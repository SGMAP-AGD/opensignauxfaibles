package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

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

func listFiles(basePath string) ([]string, error) {
	var files []string

	currentFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []string{}, err
	}

	for _, file := range currentFiles {
		if file.IsDir() {
			subPath := fmt.Sprintf("%s/%s", basePath, file.Name())
			subFiles, err := listFiles(subPath)
			if err != nil {
				return []string{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, fmt.Sprintf("%s/%s", basePath, file.Name()))
		}
	}
	return files, nil
}

// GetFileList construit la liste des fichiers à traiter
func GetFileList(basePath string, period string) (map[string][]string, error) {
	list := make(map[string][]string)
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
		l, err := ioutil.ReadDir(path)

		if err != nil {
			return nil, err
		}

		for _, f := range l {
			if match, _ := regexp.MatchString("\\.(csv|xls|xlsx)$", f.Name()); match {
				list[dir] = append(list[dir], fmt.Sprintf("%s/%s", path, f.Name()))
			}
		}
	}

	return list, nil
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
