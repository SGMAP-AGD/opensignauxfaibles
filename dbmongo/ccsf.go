package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cnf/structhash"
)

func parseCCSF(path string, CompteSiretMapping map[string]string) chan Value {
	outputChannel := make(chan Value)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

	f := map[string]int{
		"NumeroCompte":   0,
		"DateTraitement": 1,
		"Stade":          2,
		"Action":         3,
	}

	go func() {
		for {
			r, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			ccsf := CCSF{}
			ccsf.Action = r[f["Action"]]
			ccsf.Stade = r[f["Stade"]]
			ccsf.DateTraitement, err = UrssafToDate(r[f["DateTraitement"]])

			hash := fmt.Sprintf("%x", structhash.Md5(ccsf, 1))

			outputChannel <- Value{
				Value: Etablissement{
					Siret: CompteSiretMapping[r[f["NumeroCompte"]]],
					Compte: Compte{
						CCSF: map[string]CCSF{
							hash: ccsf,
						},
					},
				},
			}
		}
		close(outputChannel)
	}()
	return outputChannel
}
