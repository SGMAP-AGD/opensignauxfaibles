package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

func getCompteSiretMapping(batch *AdminBatch) (map[string]string, error) {
	compteSiretMapping := make(map[string]string)
	path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	for _, p := range path {
		file, err := os.Open(basePath + p)
		if err != nil {
			return map[string]string{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		// discard header row
		reader.Read()

		siretIndex := 3
		compteIndex := 0

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log(critical, "importCompteSiret", "Erreur à la lecture du fichier "+file.Name())
				return map[string]string{}, err
			}
			if _, err := strconv.Atoi(row[siretIndex]); err == nil && len(row[siretIndex]) == 14 {
				compteSiretMapping[row[compteIndex]] = row[siretIndex]
			}
		}
		file.Close()
	}
	return compteSiretMapping, nil
}

func importAdminUrsaff(c *gin.Context, batch *AdminBatch) {

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

func genereSeriePeriode(debut time.Time, fin time.Time) []time.Time {
	var serie []time.Time
	for fin.After(debut) {
		debut = debut.AddDate(0, 1, 0)
		serie = append(serie, debut)
	}
	return serie
}

func genereSeriePeriodeAnnuelle(debut time.Time, fin time.Time) []int {
	var serie []int
	for debut.Year() <= fin.Year() {
		serie = append(serie, debut.Year())
		debut = debut.AddDate(1, 0, 0)
	}
	return serie
}

func allErrors(slice []error, item interface{}) bool {
	for _, i := range slice {
		if i != item {
			return false
		}
	}
	return true
}
