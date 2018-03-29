package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/cnf/structhash"
)

func parseDebit(path string, CompteSiretMapping map[string]string) chan Etablissement {
	outputChannel := make(chan Etablissement)

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

	go func() {
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			if siret, ok := CompteSiretMapping[row[numeroCompteIndex]]; ok {
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

				hash := fmt.Sprintf("%x", structhash.Md5(debit, 1))
				outputChannel <- Etablissement{
					Siret: siret,
					Compte: Compte{
						Debit: map[string]Debit{
							hash: debit,
						},
					},
				}
			}
		}
		close(outputChannel)
	}()
	return outputChannel
}
