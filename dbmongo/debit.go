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

// Debit Débit – fichier Urssaf
type Debit struct {
	NumeroCompte                 string    `json:"numero_compte" bson:"numero_compte"`
	NumeroEcartNegatif           string    `json:"numero_ecart_negatif" bson:"numero_ecart_negatif"`
	DateTraitement               time.Time `json:"date_traitement" bson:"date_traitement"`
	PartOuvriere                 float64   `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale                float64   `json:"part_patronale" bson:"part_patronale"`
	NumeroHistoriqueEcartNegatif int       `json:"numero_historique" bson:"numero_historique"`
	EtatCompte                   int       `json:"etat_compte" bson:"etat_compte"`
	CodeProcedureCollective      string    `json:"code_procedure_collective" bson:"code_procedure_collective"`
	Periode                      Periode   `json:"periode" bson:"periode"`
	CodeOperationEcartNegatif    string    `json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
	CodeMotifEcartNegatif        string    `json:"code_motif_ecart_negatif" bson:"code_motif_ecart_negatif"`
	DebitSuivant                 string    `json:"debit_suivant,omitempty" bson:"debit_suivant,omitempty"`
}

func parseDebit(paths []string) chan Debit {
	outputChannel := make(chan Debit)

	go func() {
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			fields, err := reader.Read()

			numeroCompteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "num_cpte" })
			numeroEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Ecn" })
			dateTraitementIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Dt_trt_ecn" })
			partOuvriereIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PO" })
			partPatronaleIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PP" })
			numeroHistoriqueEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Hist_Ecn" })
			etatCompteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Etat_cpte" })
			codeProcedureCollectiveIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_pro_col" })
			periodeIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Periode" })
			codeOperationEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_op_ecn" })
			codeMotifEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Motif_ecn" })

			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}

				debit := Debit{}
				debit.NumeroCompte = row[numeroCompteIndex]
				debit.NumeroEcartNegatif = row[numeroEcartNegatifIndex]
				debit.DateTraitement, err = UrssafToDate(row[dateTraitementIndex])
				debit.PartOuvriere, err = strconv.ParseFloat(row[partOuvriereIndex], 64)
				debit.PartOuvriere = debit.PartOuvriere / 100
				debit.PartPatronale, err = strconv.ParseFloat(row[partPatronaleIndex], 64)
				debit.PartPatronale = debit.PartPatronale / 100
				debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[numeroHistoriqueEcartNegatifIndex])
				debit.EtatCompte, err = strconv.Atoi(row[etatCompteIndex])
				debit.CodeProcedureCollective = row[codeProcedureCollectiveIndex]
				debit.Periode, err = UrssafToPeriod(row[periodeIndex])
				debit.CodeOperationEcartNegatif = row[codeOperationEcartNegatifIndex]
				debit.CodeMotifEcartNegatif = row[codeMotifEcartNegatifIndex]

				outputChannel <- debit
			}
			file.Close()
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importDebit(c *gin.Context) {
	insertWorker := c.Keys["insertEtablissement"].(chan ValueEtablissement)

	batch := c.Params.ByName("batch")

	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	dataSource := files["debit"]
	mapping := getCompteSiretMapping(files["admin_urssaf"])

	for debit := range parseDebit(dataSource) {
		if siret, ok := mapping[debit.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(debit, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch: Batch{
							// Compact: map[string]bool{
							// 	"status": false,
							// },
							Debit: map[string]Debit{
								hash: debit,
							}}}}}
			insertWorker <- value
		}
	}
}
