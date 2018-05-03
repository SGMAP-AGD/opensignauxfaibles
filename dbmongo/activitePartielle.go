package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/tealeg/xlsx"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `json:"id_demande" bson:"id_demande"`
	EffectifEntreprise int       `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           int       `json:"effectif" bson:"effectif"`
	DateStatut         time.Time `json:"date_statut" bson:"date_statut"`
	TxPC               float64   `json:"tx_pc" bson:"tx_pc"`
	TxPCUnedicDares    float64   `json:"tx_pc_unedic_dares" bson:"tx_pc_unedic_dares"`
	TxPCEtatDares      float64   `json:"tx_pc_etat_dares" bson:"tx_pc_etat_dares"`
	Periode            Periode   `json:"periode" bson:"periode"`
	HTA                float64   `json:"hta" bson:"hta"`
	MTA                float64   `json:"mta" bson:"mta"`
	EffectifAutorise   int       `json:"effectif_autorise" bson:"effectif_autorise"`
	ProdHTAEffectif    float64   `json:"prod_hta_effectif" bson:"prod_hta_effectif"`
	MotifRecoursSE     int       `json:"motif_recours_se" bson:"motif_recours_se"`
	Perimetre          int       `json:"perimetre" bson:"perimetre"`
	RecoursAnterieur   int       `json:"recours_anterieur" bson:"recours_anterieur"`
	AvisCE             int       `json:"avis_ce" bson:"avis_ce"`
	HeureConsommee     float64   `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    float64   `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   int       `json:"effectif_consomme" bson:"effectif_consomme"`
}

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `json:"id_conso" bson:"id_conso"`
	HeureConsommee float64   `json:"heure_consomme" bson:"heure_consomme"`
	Montant        float64   `json:"montant" bson:"montant"`
	Effectif       int       `json:"effectif" bson:"effectif"`
	Periode        time.Time `json:"periode" bson:"periode"`
}

func parseAPDemande(path string, batch string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	go func() {
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows[3:] {
				apdemande := APDemande{}
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
					Key: row.Cells[3].Value,
					Batch: map[string]Batch{
						batch: Batch{
							APDemande: map[string]APDemande{
								hash: apdemande,
							},
						},
					},
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}

func parseAPConsommation(path string, batch string) chan Etablissement {
	outputChannel := make(chan Etablissement)

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("Error", err)
	}

	go func() {
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows[3:] {
				apconsommation := APConso{}
				apconsommation.ID = row.Cells[1].Value
				apconsommation.Periode, err = excelToTime(row.Cells[15].Value)
				apconsommation.HeureConsommee, err = strconv.ParseFloat(row.Cells[16].Value, 64)
				apconsommation.Montant, err = strconv.ParseFloat(row.Cells[17].Value, 64)
				apconsommation.Effectif, err = strconv.Atoi(row.Cells[18].Value)

				if err != nil {
					fmt.Println(err)
				}

				hash := fmt.Sprintf("%x", structhash.Md5(apconsommation, 1))
				outputChannel <- Etablissement{
					Key: row.Cells[3].Value,
					Batch: map[string]Batch{
						batch: Batch{
							APConso: map[string]APConso{
								hash: apconsommation,
							},
						},
					},
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}
