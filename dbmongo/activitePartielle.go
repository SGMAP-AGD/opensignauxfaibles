package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/tealeg/xlsx"
)

func excelToTime(excel string) (time.Time, error) {
	excelInt, err := strconv.ParseInt(excel, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}
	return time.Unix((excelInt-25569)*3600*24, 0), nil
}
func parseActivitePartielleDemande(path string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	go func() {
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows[2:] {
				apdemande := APDemande{}
				apdemande.ID = row.Cells[2].Value
				apdemande.EffectifEntreprise, _ = strconv.Atoi(row.Cells[14].Value)
				apdemande.Effectif, _ = strconv.Atoi(row.Cells[15].Value)
				apdemande.DateStatut, _ = excelToTime(row.Cells[16].Value)
				apdemande.TxPC, _ = strconv.ParseFloat(row.Cells[17].Value, 64)
				apdemande.TxPCUnedicDares, _ = strconv.ParseFloat(row.Cells[18].Value, 64)
				apdemande.TxPCEtatDares, _ = strconv.ParseFloat(row.Cells[19].Value, 64)
				periodStart, _ := excelToTime(row.Cells[20].Value)
				periodEnd, _ := excelToTime(row.Cells[21].Value)
				apdemande.Periode = Periode{
					Start: periodStart,
					End:   periodEnd,
				}
				apdemande.HTA, _ = strconv.ParseFloat(row.Cells[22].Value, 64)
				apdemande.MTA, _ = strconv.ParseFloat(row.Cells[23].Value, 64)
				apdemande.EffectifAutorise, _ = strconv.Atoi(row.Cells[24].Value)
				apdemande.ProdHTAEffectif, _ = strconv.ParseFloat(row.Cells[25].Value, 64)
				apdemande.MotifRecoursSE, _ = strconv.Atoi(row.Cells[26].Value)
				apdemande.Perimetre, _ = strconv.Atoi(row.Cells[27].Value)
				apdemande.RecoursAnterieur, _ = strconv.Atoi(row.Cells[28].Value)
				apdemande.AvisCE, _ = strconv.Atoi(row.Cells[29].Value)
				apdemande.HeureConsommee, _ = strconv.ParseFloat(row.Cells[30].Value, 64)
				apdemande.MontantConsomme, _ = strconv.ParseFloat(row.Cells[31].Value, 64)
				apdemande.EffectifConsomme, _ = strconv.Atoi(row.Cells[32].Value)

				hash := fmt.Sprintf("%x", structhash.Md5(apdemande, 1))

				outputChannel <- Etablissement{
					Siret: row.Cells[3].Value,
					ActivitePartielle: ActivitePartielle{
						Demande: map[string]APDemande{
							hash: apdemande,
						},
					},
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}

func parseActivitePartielleConsommation(path string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	go func() {
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows[3:] {
				apconsommation := APConsommation{}
				apconsommation.ID = row.Cells[1].Value
				apconsommation.Date, err = excelToTime(row.Cells[15].Value)
				apconsommation.HeureConsommee, err = strconv.ParseFloat(row.Cells[16].Value, 64)
				apconsommation.Montant, err = strconv.ParseFloat(row.Cells[17].Value, 64)
				apconsommation.Effectif, err = strconv.Atoi(row.Cells[18].Value)

				if err != nil {
					fmt.Println(err)
				}

				hash := fmt.Sprintf("%x", structhash.Md5(apconsommation, 1))
				outputChannel <- Etablissement{
					Siret: row.Cells[2].Value,
					ActivitePartielle: ActivitePartielle{
						Consommation: map[string]APConsommation{
							hash: apconsommation,
						},
					},
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}
