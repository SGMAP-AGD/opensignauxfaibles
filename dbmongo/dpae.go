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

// DPAE Déclaration préalabre à l'embauche
type DPAE struct {
	Siret    string    `json:"-" bson:"-"`
	Date     time.Time `json:"date" bson:"date"`
	CDI      float64   `json:"cdi" bson:"cdi"`
	CDDLong  float64   `json:"cdd_long" bson:"cdd_long"`
	CDDCourt float64   `json:"cdd_court" bson:"cdd_court"`
}

func parseDPAE(path string) chan *DPAE {
	outputChannel := make(chan *DPAE)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

	go func() {
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			date, err := time.Parse("20060102", row[1]+row[2]+"01")

			dpae := DPAE{
				Siret: row[0],
				Date:  date,
			}
			dpae.CDI, _ = strconv.ParseFloat(row[3], 64)
			dpae.CDDLong, _ = strconv.ParseFloat(row[4], 64)
			dpae.CDDCourt, _ = strconv.ParseFloat(row[5], 64)

			if err == nil {
				outputChannel <- &dpae

			}
		}
		file.Close()
		close(outputChannel)
	}()
	return outputChannel
}

func importDPAE(c *gin.Context) {
	insertWorker, _ := c.Keys["insertEtablissement"].(chan ValueEtablissement)
	batch := c.Params.ByName("batch")
	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	dpaes := files["dpae"]

	for _, dpaeFile := range dpaes {
		for dpae := range parseDPAE(dpaeFile) {
			hash := fmt.Sprintf("%x", structhash.Md5(dpae, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: dpae.Siret,
					Batch: map[string]Batch{
						batch: Batch{
							DPAE: map[string]*DPAE{
								hash: dpae,
							}}}}}
			insertWorker <- value
		}
	}
}
