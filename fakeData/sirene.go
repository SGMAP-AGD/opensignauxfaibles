package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func newWords(wordList []string) {

}
func readAndRandomSirene(fileName string, outputFileName string, mapping map[string]string) error {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','

	title, err := reader.Read()
	outputRow := strings.Join(title, ",") + "\n"
	outputLength := len(title)

	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return err
	}

	var wordBase = make(map[string]struct{})

	for i := 0; i < 100000; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}
		words := strings.Split(row[2], " ")
		for _, word := range words {
			wordBase[word] = struct{}{}
		}
	}

	var wordList []string
	for key := range wordBase {
		wordList = append(wordList, key)
	}

	wordLength := float64(len(wordList))
	length := int(rand.Float32()*5) + 1

	// sirene.APE = row[42]
	// sirene.NatureActivite = row[52]
	// sirene.ActiviteSaisoniere = row[55]
	// sirene.ModaliteActivite = row[56]
	// sirene.Productif = row[57]
	// sirene.NatureJuridique = row[71]
	// sirene.Categorie = row[82]
	// sirene.Creation, _ = time.Parse("20060102", row[50])
	// sirene.IndiceMonoactivite, _ = strconv.Atoi(row[85])
	// sirene.TrancheCA, _ = strconv.Atoi(row[89])
	// sirene.Sigle = row[61]
	// sirene.DebutActivite, _ = time.Parse("20060102", row[51])
	// sirene.Longitude, _ = strconv.ParseFloat(row[100], 64)
	// sirene.Lattitude, _ = strconv.ParseFloat(row[101], 64)
	// sirene.Adresse = [7]string{row[2], row[3], row[4], row[5], row[6], row[7], row[8]}

	for _, v := range mapping {
		if len(v) == 14 {
			output := make([]string, outputLength)
			// siren
			output[0] = v[0:9]
			// nic
			output[1] = v[9:13]
			// nic siege
			output[65] = v[9:13]

			output[16] = "21"
			output[19] = "RUE"
			output[20] = "21000"
			output[23] = "Bourgogne Franche-Comté"
			output[24] = "Côte d'or"
			output[28] = "Dijon"
			output[42] = "2120Z"
			output[71] = "SA"
			output[81] = "2018-01-01"
			output[51] = "2018-01-01"
			output[100] = "5.051709"
			output[101] = "47.315471"
			output[2] = "21 BOULEVARD VOLTAIRE"
			output[4] = "21000 DIJON"

			var sentence []string
			for i := 0; i < length; i++ {
				sentence = append(sentence, wordList[int(rand.Float64()*wordLength)])
			}
			// raison sociale
			output[2] = strings.Join(sentence, " ")

			outputRow := strings.Join(output, ",") + "\n"
			_, err = outputFile.WriteString(outputRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
