package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Pour déformer des fichiers réels pour créer un dataset consistent avec anonymisation des entreprises

func main() {
	rand.Seed(time.Now().UnixNano())

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration:" + err.Error())
	}

	// Traitement des comptes
	prefixOutput := viper.GetString("prefixOutput")
	comptes := viper.GetString("comptes")

	outputCompte := outputFileName(prefixOutput, comptes)
	mapping, err := readAndRandomComptes(comptes, outputCompte)

	fmt.Println(mapping)
	// // Traitement des débits
	// debits := viper.GetString("debits")
	// outputDebits := outputFileName(viper.GetString("prefixOutput"), debits)
	// fmt.Print("Fake débits: ")
	// err = readAndRandomDebits(debits, outputDebits)
	// if err != nil {
	// 	fmt.Println("Fail : " + err.Error())
	// 	fmt.Println("Interruption.")
	// } else {
	// 	fmt.Println("OK -> " + outputDebits)
	// }

	// // Traitement des cotisations
	// cotisations := viper.GetString("cotisations")
	// outputCotisations := outputFileName(viper.GetString("prefixOutput"), cotisations)
	// fmt.Print("Fake cotisations: ")
	// err = readAndRandomCotisations(cotisations, outputCotisations)
	// if err != nil {
	// 	fmt.Println("Fail : " + err.Error())
	// 	fmt.Println("Interruption.")
	// } else {
	// 	fmt.Println("OK -> " + outputCotisations)
	// }
}

func outputFileName(prefixOutput string, fileName string) string {
	path := strings.Split(fileName, "/")
	path[len(path)-1] = prefixOutput + path[len(path)-1]

	return strings.Join(path, "/")
}

const letterBytes = "0123456789"

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
