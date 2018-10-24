package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func readAndRandomDebits(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	coef := make(map[string]float64)

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	row, err := reader.Read()
	fmt.Println("\"" + strings.Join(row, "\";\"") + "\"")

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		partOuvriere, _ := strconv.ParseFloat(row[4], 64)
		partPatronale, _ := strconv.ParseFloat(row[5], 64)
		partPenalite, _ := strconv.ParseFloat(row[15], 64)

		if c, ok := coef[row[0]]; ok {
			row[4] = strconv.Itoa(int(partOuvriere * c))
			row[5] = strconv.Itoa(int(partPatronale * c))
			row[15] = strconv.Itoa(int(partPenalite * c))
		} else {
			coef[row[0]] = rand.Float64() * rand.Float64() / 150
			row[4] = strconv.Itoa(int(partOuvriere * coef[row[0]]))
			row[5] = strconv.Itoa(int(partPatronale * coef[row[0]]))
			row[15] = strconv.Itoa(int(partPenalite * coef[row[0]]))
		}
		fmt.Println("\"" + strings.Join(row, "\";\"") + "\"")
	}
	return nil
}
