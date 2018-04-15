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

func parseCotisation(path string, CompteSiretMapping map[string]string) chan Value {
	outputChannel := make(chan Value)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

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
		cotisation := Cotisation{}

		for {

			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			if siret, ok := CompteSiretMapping[row[field["NumeroCompte"]]]; ok {

				cotisation = Cotisation{}
				cotisation.NumeroCompte = row[field["NumeroCompte"]]
				cotisation.Periode, err = UrssafToPeriod(row[field["Periode"]])
				cotisation.PeriodeDebit = row[field["PeriodeDebit"]]
				cotisation.Recouvrement, err = strconv.ParseFloat(strings.Replace(row[field["Recouvrement"]], ",", ".", -1), 64)
				cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
				cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
				cotisation.Ecriture = row[field["Ecriture"]]

				hash := fmt.Sprintf("%x", structhash.Md5(cotisation, 1))
				outputChannel <- Value{
					Value: Etablissement{
						Siret: siret,
						Compte: Compte{
							Cotisation: map[string]Cotisation{
								hash: cotisation,
							},
						},
					},
				}
			}
		}

		close(outputChannel)
	}()

	return outputChannel
}
