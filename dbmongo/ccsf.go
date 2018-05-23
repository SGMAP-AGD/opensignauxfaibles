package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CCSF information urssaf ccsf
type CCSF struct {
	NumeroCompte   string    `json:"-" bson:"-"`
	DateTraitement time.Time `json:"date_traitement" bson:"date_traitement"`
	Stade          string    `json:"stade" bson:"stade"`
	Action         string    `json:"action" json:"action"`
}

func parseCCSF(path string) chan CCSF {
	outputChannel := make(chan CCSF)

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
			ccsf.NumeroCompte = r[f["NumeroCompte"]]
			outputChannel <- ccsf
		}
		close(outputChannel)
		file.Close()
	}()
	return outputChannel
}

func importCCSF(c *gin.Context) {
	insertWorker := c.Keys["DBW"].(chan Value)

	batch := c.Params.ByName("batch")

	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	dataSource := files["ccsf"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	for _, data := range dataSource {
		for ccsf := range parseCCSF(data) {
			if siret, ok := mapping[ccsf.NumeroCompte]; ok {
				hash := fmt.Sprintf("%x", structhash.Md5(ccsf, 1))

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
										CCSF: map[string]CCSF{
											hash: ccsf,
										}}}}}}}
				insertWorker <- value
			}
		}
	}
}
