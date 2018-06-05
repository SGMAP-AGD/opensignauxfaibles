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

func sliceIndexArray(target []string, data []string) []int {
	l := len(target)
	var idx []int

	for a := 1; a <= l; a++ {
		idx = append(idx, -1)
	}

	for idxTarget, t := range target {
		for idxData, d := range data {
			if d == t {
				idx[idxTarget] = idxData
			}
		}
	}

	return idx
}

type stringSlice []string

func (arr stringSlice) contains(str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// batchToTime calcule la date de référence à partir de la référence de batch
func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
}

func timeToBatch(date time.Time) string {
	return date.Format("0602")
}

func batchList(first string, last string) ([]string, error) {
	return nil, nil
}
