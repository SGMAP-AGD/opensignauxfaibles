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
)

func parseDelais(path string, CompteSiretMapping map[string]string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

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
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			delais := Delais{}
			delais.NumeroCompte = row[field["NumeroCompte"]]
			delais.NumeroContentieux = row[field["NumeroContentieux"]]
			delais.DateCreation, err = time.Parse("2006-01-02", row[field["DateCreation"]])
			delais.DateEcheanche, err = time.Parse("2006-01-02", row[field["DateEcheanche"]])
			delais.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
			delais.Denomination = row[field["Denomination"]]
			delais.Indic6m = row[field["Indic6m"]]
			delais.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
			delais.MontantEcheancier, err = strconv.ParseFloat(row[field["MontantEcheancier"]], 64)
			delais.NumeroStructure = row[field["NumeroStructure"]]
			delais.Stade = row[field["Stade"]]
			delais.Action = row[field["Action"]]

			hash := fmt.Sprintf("%x", structhash.Md5(delais, 1))

			if siret, ok := CompteSiretMapping[row[field["NumeroCompte"]]]; ok {
				outputChannel <- Etablissement{
					Siret: siret,
					Compte: Compte{
						Delais: map[string]Delais{
							hash: delais,
						},
					},
				}
			}
		}
		close(outputChannel)

	}()

	return outputChannel
}
