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
)

func parseAltares(path string) chan Etablissement {
	outputChannel := make(chan Etablissement)

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
				DateEffet:     dateEffet,
				DateParution:  dateParution,
				CodeJournal:   row[codeJournalIndex],
				CodeEvenement: row[codeEvenementIndex],
			}
			hash := fmt.Sprintf("%x", structhash.Md5(altares, 1))
			if err == nil {
				outputChannel <- Etablissement{
					Siret: row[siretIndex],
					Altares: map[string]Altares{
						hash: altares,
					},
				}
			}
		}
		close(outputChannel)
	}()
	return outputChannel
}
