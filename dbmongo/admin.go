package main

import (
	"fmt"

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
