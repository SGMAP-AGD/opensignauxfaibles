package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"

	"github.com/spf13/viper"
)

// Pour déformer des fichiers réels pour créer un dataset consistent avec anonymisation des entreprises

type randomizer func(string, string, map[string]string) error

// run execute une fonction de randomisation
func run(name string, handler randomizer, mapping map[string]string) error {
	file := viper.GetString(name)
	outputFile := outputFileName(viper.GetString("prefixOutput"), file)
	fmt.Print("Fake " + name + ": ")
	err := handler(file, outputFile, mapping)
	if err != nil {
		fmt.Println("Fail : " + err.Error())
		fmt.Println("Interruption.")
		return fmt.Errorf("Interruption")
	}
	fmt.Println("OK -> " + outputFile)
	return nil
}

func main() {

	checkPtr := flag.String("check", "", "Check configuration.")
	randomPtr := flag.String("random", "", "Randomize datasource.")
	uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
	flag.Parse()

	fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *checkPtr, *randomPtr, *uniquePtr)

	// viper.SetConfigName("config")
	// viper.SetConfigType("toml")
	// viper.AddConfigPath(".")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic("Erreur à la lecture de la configuration:" + err.Error())
	// }

	// // Génération des numéros de comptes
	// outputCompte := outputFileName(viper.GetString("prefixOutput"), viper.GetString("comptes"))
	// fmt.Print("Fake comptes: ")
	// mapping, err := readAndRandomComptes(viper.GetString("comptes"), outputCompte)
	// if err != nil {
	// 	fmt.Println("Fail : " + err.Error())
	// 	fmt.Println("Interruption.")
	// } else {
	// 	fmt.Println("OK -> " + outputCompte)
	// }

	// // Traitement des délais
	// // err = run("delais", readAndRandomDelais, mapping)
	// // if err != nil {
	// // 	panic(err)
	// // }

	// // Traitement sirene
	// err = run("sirene", readAndRandomSirene, mapping)
	// if err != nil {
	// 	panic(err)
	// }

	// // err = run("debits", readAndRandomDebits, mapping)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // // Traitement des débits
	// // debits := viper.GetString("debits")
	// // outputDebits := outputFileName(viper.GetString("prefixOutput"), debits)
	// // fmt.Print("Fake débits: ")
	// // err = readAndRandomDebits(debits, outputDebits, mapping)
	// // if err != nil {
	// // 	fmt.Println("Fail : " + err.Error())
	// // 	fmt.Println("Interruption.")
	// // } else {
	// // 	fmt.Println("OK -> " + outputDebits)
	// // }

	// // // Traitement altares
	// // altares := viper.GetString("altares")
	// // outputAltares := outputFileName(prefixOutput, altares)
	// // fmt.Print("Fake altares: ")
	// // err = readAndRandomAltares(altares, outputAltares, mapping)
	// // if err != nil {
	// // 	fmt.Println("Fail : " + err.Error())
	// // 	panic("Interruption.")
	// // } else {
	// // 	fmt.Println("OK -> " + outputAltares)
	// // }
	// // // Traitement des débits
	// // debits := viper.GetString("debits")
	// // outputDebits := outputFileName(viper.GetString("prefixOutput"), debits)
	// // fmt.Print("Fake débits: ")
	// // err = readAndRandomDebits(debits, outputDebits, mapping)
	// // if err != nil {
	// // 	fmt.Println("Fail : " + err.Error())
	// // 	fmt.Println("Interruption.")
	// // } else {
	// // 	fmt.Println("OK -> " + outputDebits)
	// // }

	// // // Traitement des cotisations
	// // cotisations := viper.GetString("cotisations")
	// // outputCotisations := outputFileName(viper.GetString("prefixOutput"), cotisations)
	// // fmt.Print("Fake cotisations: ")
	// // err = readAndRandomCotisations(cotisations, outputCotisations, mapping)
	// // if err != nil {
	// // 	fmt.Println("Fail : " + err.Error())
	// // 	fmt.Println("Interruption.")
	// // } else {
	// // 	fmt.Println("OK -> " + outputCotisations)
	// // }
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
