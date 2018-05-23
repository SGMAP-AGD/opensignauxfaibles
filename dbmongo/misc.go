package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func excelToTime(excel string) (time.Time, error) {
	excelInt, err := strconv.ParseInt(excel, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}
	return time.Unix((excelInt-25569)*3600*24, 0), nil
}

var (
	trace    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	logerror *log.Logger
)

// InitLogger initialise les variables permettant l'écriture des messages de log
func InitLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	logerror = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func getCompteSiretMapping(path []string) map[string]string {
	compteSiretMapping := make(map[string]string)

	for _, p := range path {
		file, _ := os.Open(p)

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		// discard header row
		reader.Read()

		siretIndex := 3
		compteIndex := 0

		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			if _, err := strconv.Atoi(row[siretIndex]); err == nil && len(row[siretIndex]) == 14 {
				compteSiretMapping[row[compteIndex]] = row[siretIndex]
			}
		}
		file.Close()
	}
	return compteSiretMapping
}
