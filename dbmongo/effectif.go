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

// Effectif Urssaf
type Effectif struct {
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
}

// ParseEffectifPeriod Transforme un tableau de périodes telles qu'écrites dans l'entête du tableau d'effectif urssaf en date de début
func parseEffectifPeriod(effectifPeriods []string) ([]time.Time, error) {
	periods := []time.Time{}
	for _, period := range effectifPeriods {
		urssaf := period[3:9]
		date, _ := UrssafToPeriod(urssaf)
		periods = append(periods, date.Start)
	}

	return periods, nil
}

func getCompteSiretMapping(path string) map[string]string {
	compteSiretMapping := make(map[string]string)

	file, _ := os.Open(path)

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	fields, _ := reader.Read()

	siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "SIRET" })
	compteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "compte" })

	for {
		row, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		compteSiretMapping[row[compteIndex]] = row[siretIndex]
	}

	return compteSiretMapping
}

func parseEffectif(path string) chan Value {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	fields, err := reader.Read()

	if err != nil {
		fmt.Println("Error", err)
	}

	siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "SIRET" })
	compteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "compte" })

	boundaryIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "ape_ins" })
	periods, err := parseEffectifPeriod(fields[0:boundaryIndex])

	if err != nil {
		log.Panic("Aborting: could not read a period:", err)
	}
	outputChannel := make(chan Value)

	go func() {
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			i := 0
			effectif := map[string]Effectif{}

			for i < boundaryIndex {
				e, _ := strconv.Atoi(row[i])
				if e > 0 {
					eff := Effectif{
						NumeroCompte: row[compteIndex],
						Periode:      periods[i],
						Effectif:     e}
					hash := fmt.Sprintf("%x", structhash.Md5(eff, 1))
					effectif[hash] = eff
				}
				i++
			}

			if len(row[siretIndex]) == 14 {
				outputChannel <- Value{
					Value: Etablissement{
						Siret: row[siretIndex],
						Compte: Compte{
							Effectif: effectif,
						},
					},
				}

			}
		}

		close(outputChannel)
	}()

	return outputChannel
}
