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

// Altares Extrait du récapitulatif altarès
type Altares struct {
	DateEffet     time.Time `json:"date_effet" bson:"date_effet"`
	DateParution  time.Time `json:"date_parution" bson:"date_parution"`
	CodeJournal   string    `json:"code_journal" bson:"code_journal"`
	CodeEvenement string    `json:"code_evenement" bson:"code_evenement"`
	Siret         string    `json:"-" bson:"-"`
}

func parseAltares(path string, batch string) chan Altares {
	outputChannel := make(chan Altares)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','
	reader.LazyQuotes = true
	fields, err := reader.Read()

	dateEffetIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date d'effet" })
	dateParutionIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date parution" })
	codeJournalIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code du journal" })
	codeEvenementIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code de la nature de l'événement" })
	siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Siret" })

	go func() {
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			dateEffet, err := time.Parse("2006-01-02", row[dateEffetIndex])
			dateParution, _ := time.Parse("2006-01-02", row[dateParutionIndex])

			altares := Altares{
				Siret:         row[siretIndex],
				DateEffet:     dateEffet,
				DateParution:  dateParution,
				CodeJournal:   row[codeJournalIndex],
				CodeEvenement: row[codeEvenementIndex],
			}
			if err == nil {
				outputChannel <- altares

			}
		}
		file.Close()
		close(outputChannel)
	}()
	return outputChannel
}

func importAltares(c *gin.Context) {
	insertWorker, _ := c.Keys["insertEtablissement"].(chan ValueEtablissement)
	batch := c.Params.ByName("batch")
	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	altares := files["altares"][0]

	for altares := range parseAltares(altares, batch) {
		hash := fmt.Sprintf("%x", structhash.Md5(altares, 1))

		value := ValueEtablissement{
			Value: Etablissement{
				Siret: altares.Siret,
				Batch: map[string]Batch{
					batch: Batch{
						// Compact: map[string]bool{
						// 	"status": false,
						// },
						Altares: map[string]Altares{
							hash: altares,
						}}}}}
		insertWorker <- value
	}
}
