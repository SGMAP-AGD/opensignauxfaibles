package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// AdminID Collection key
type AdminID struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}

func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
}

// Types description des types de fichiers pris en charge
type Types []struct {
	Type    string `json:"type" bson:"type"`
	Libelle string `json:"text" bson:"text"`
	Filter  string `json:"filter" bson:"filter"`
}

func listTypes(c *gin.Context) {
	c.JSON(200, Types{
		{"admin_urssaf", "Siret/Compte URSSAF", "Liste comptes"},
		{"apconso", "Consommation Activité Partielle", "conso"},
		{"bdf", "Ratios Banque de France", "bdf"},
		{"cotisation", "Cotisations URSSAF", "cotisation"},
		{"delai", "Délais URSSAF", "delais|Délais"},
		{"dpae", "Déclaration Préalable à l'embauche", "DPAE"},
		{"interim", "Base Interim", "interim"},
		{"altares", "Base Altarès", "ALTARES"},
		{"apdemande", "Demande Activité Partielle", "dde"},
		{"ccsf", "Stock CCSF à date", "ccsf"},
		{"debit", "Débits URSSAF", "debit"},
		{"dmmo", "Déclaration Mouvement de Main d'Œuvre", "dmmo"},
		{"effectif", "Emplois URSSAF", "Emploi"},
		{"sirene", "Base GéoSirene", "sirene"},
		{"diane", "Diane", "diane"},
	})
}

func cloneDB(c *gin.Context) {
	from := viper.GetString("DB")
	to := c.Params.ByName("to")

	var result interface{}
	declareDatabaseCopy(db.DB, from, to)
	err := db.DB.Run(bson.M{"eval": "copyDatabase()"}, result)
	if err != nil {
		c.JSON(500, err)
		fmt.Println(err)
		return
	}
	removeDatabaseCopy(db.DB)
	c.JSON(200, err)
}

func addFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(500, "Error occured: "+err.Error())
		return
	}

	source, err := header.Open()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	defer source.Close()

	destination, err := os.Create(viper.GetString("APP_DATA") + "/" + header.Filename)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)

	type result struct {
		Bytes int64 `json:"bytes,omitempty"`
		Error error `json:"error,omitempty"`
	}

	if err != nil {
		c.JSON(500, result{nBytes, err})
		return
	}

	basePath := viper.GetString("APP_DATA")
	files, err := listFiles(basePath)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	mainMessageChannel <- socketMessage{
		Files: files,
	}

	c.JSON(200, result{nBytes, err})
}
