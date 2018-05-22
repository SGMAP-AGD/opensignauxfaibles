package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Cotisation Cotisation – fichier Urssaf
type Cotisation struct {
	NumeroCompte string  `json:"numero_compte" bson:"numero_compte"`
	PeriodeDebit string  `json:"periode_debit" bson:"periode_debit"`
	Periode      Periode `json:"period" bson:"periode"`
	Recouvrement float64 `json:"recouvrement" bson:"recouvrement"`
	Encaisse     float64 `json:"encaisse" bson:"encaisse"`
	Du           float64 `json:"du" bson:"du"`
	Ecriture     string  `json:"ecriture" bson:"ecriture"`
}

func parseCotisation(paths []string) chan Cotisation {
	outputChannel := make(chan Cotisation)

	field := map[string]int{
		"NumeroCompte": 0,
		"PeriodeDebit": 1,
		"Periode":      4,
		"Recouvrement": 2,
		"Encaisse":     3,
		"Du":           5,
		"Ecriture":     6,
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

			cotisation := Cotisation{}

			for {

				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}

				cotisation = Cotisation{}
				cotisation.NumeroCompte = row[field["NumeroCompte"]]
				cotisation.Periode, err = UrssafToPeriod(row[field["Periode"]])
				cotisation.PeriodeDebit = row[field["PeriodeDebit"]]
				cotisation.Recouvrement, err = strconv.ParseFloat(strings.Replace(row[field["Recouvrement"]], ",", ".", -1), 64)
				cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
				cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
				cotisation.Ecriture = row[field["Ecriture"]]

				outputChannel <- cotisation
			}
			file.Close()
		}

		close(outputChannel)
	}()
	return outputChannel
}

func importCotisation(c *gin.Context) {
	insertWorker := c.Keys["DBW"].(chan Value)

	batch := c.Params.ByName("batch")

	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	dataSource := files["cotisation"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	for cotisation := range parseCotisation(dataSource) {
		if siret, ok := mapping[cotisation.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(cotisation, 1))

			value := Value{
				Value: Entreprise{
					Siren: siret[0:9],
					Etablissement: map[string]Etablissement{
						siret: Etablissement{
							Siret: siret,
							Batch: map[string]Batch{
								batch: Batch{
									Compact: map[string]bool{
										"status": false,
									},
									Cotisation: map[string]Cotisation{
										hash: cotisation,
									}}}}}}}
			insertWorker <- value
		}
	}
}
