package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

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

func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
}

func listTypes(c *gin.Context) {
	c.JSON(200, []struct {
		Type    string `json:"type" bson:"type"`
		Libelle string `json:"text" bson:"text"`
		Filter  string `json:"filter" bson:"filter"`
	}{
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
	})
}

func cloneDB(c *gin.Context) {
	db := c.Keys["db"].(*mgo.Database)

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
func listBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var batch []AdminBatch
	db.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batch)
	c.JSON(200, batch)
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
