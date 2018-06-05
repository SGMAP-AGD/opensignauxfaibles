package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Delai tuple fichier ursaff
type Delai struct {
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheanche     time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	NumeroStructure   string    `json:"numero_structure" bson:"numero_structure"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

func parseDelai(paths []string) chan *Delai {
	outputChannel := make(chan *Delai)

	field := map[string]int{
		"NumeroCompte":      0,
		"NumeroContentieux": 1,
		"DateCreation":      2,
		"DateEcheanche":     3,
		"DureeDelai":        4,
		"Denomination":      5,
		"Indic6m":           6,
		"AnneeCreation":     7,
		"MontantEcheancier": 8,
		"NumeroStructure":   9,
		"Stade":             10,
		"Action":            11,
	}

	go func() {
		for _, path := range paths {

			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()

			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}

				delai := Delai{}
				delai.NumeroCompte = row[field["NumeroCompte"]]
				delai.NumeroContentieux = row[field["NumeroContentieux"]]
				delai.DateCreation, err = time.Parse("2006-01-02", row[field["DateCreation"]])
				delai.DateEcheanche, err = time.Parse("2006-01-02", row[field["DateEcheanche"]])
				delai.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
				delai.Denomination = row[field["Denomination"]]
				delai.Indic6m = row[field["Indic6m"]]
				delai.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
				delai.MontantEcheancier, err = strconv.ParseFloat(row[field["MontantEcheancier"]], 64)
				delai.NumeroStructure = row[field["NumeroStructure"]]
				delai.Stade = row[field["Stade"]]
				delai.Action = row[field["Action"]]

				outputChannel <- &delai
			}
			file.Close()
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importDelai(c *gin.Context) {
	insertWorker := c.Keys["insertEtablissement"].(chan *ValueEtablissement)

	batch := c.Params.ByName("batch")

	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	dataSource := files["delai"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	for delai := range parseDelai(dataSource) {
		if siret, ok := mapping[delai.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(delai, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch: Batch{
							Delai: map[string]*Delai{
								hash: delai,
							}}}}}
			insertWorker <- &value
		}
	}
	insertWorker <- &ValueEtablissement{}

}
