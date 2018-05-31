package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `json:"id_demande" bson:"id_demande"`
	Siret              string    `json:"-" bson:"-"`
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
	Siret          string    `json:"-" bson:"-"`
	HeureConsommee float64   `json:"heure_consomme" bson:"heure_consomme"`
	Montant        float64   `json:"montant" bson:"montant"`
	Effectif       int       `json:"effectif" bson:"effectif"`
	Periode        time.Time `json:"periode" bson:"periode"`
}

func parseAPDemande(path string) chan *APDemande {
	outputChannel := make(chan *APDemande)

	go func() {
		xlFile, err := xlsx.OpenFile(path)
		if err != nil {
			fmt.Println("parseAPDemande: Avorté " + err.Error())
			close(outputChannel)
			return
		}

		sheet := xlFile.Sheets[0]
		f := make(map[string]int)
		for idx, cell := range sheet.Rows[0].Cells {
			f[cell.Value] = idx
		}

		fields := []string{
			"ID_DA",
			"ETAB_SIRET",
			"EFF_ENT",
			"EFF_ETAB",
			"DATE_STATUT",
			// "TX_PC",
			// "TX_PC_UNEDIC_DARES",
			// "TX_PC_ETAT_DARES",
			"DATE_DEB",
			"DATE_FIN",
			"HTA",
			// "MTA",
			"EFF_AUTO",
			// "PROD_HTA_EFF",
			"MOTIF_RECOURS_SE",
			//"PERIMETRE_AP",
			// "RECOURS_ANTERIEUR",
			// "AVIS_CE",
			"S_HEURE_CONSOM_TOT",
			// "S_MONTANT_CONSOM_TOT",
			"S_EFF_CONSOM_TOT",
		}
		minLength := 0
		for _, field := range fields {
			if i, err := f[field]; err {
				minLength = max(minLength, i)
			} else {
				fmt.Println("parseAPDemande: Avorté, " + field + " non trouvé.")
				close(outputChannel)
				return
			}
		}

		// FIXME: établir une meilleure méthode pour tester la validité de la ligne
		for _, row := range sheet.Rows[3:] {
			if len(row.Cells) >= minLength {
				apdemande := APDemande{}
				apdemande.ID = row.Cells[f["ID_DA"]].Value
				apdemande.Siret = row.Cells[f["ETAB_SIRET"]].Value
				apdemande.EffectifEntreprise, _ = strconv.Atoi(row.Cells[f["EFF_ENT"]].Value)
				apdemande.Effectif, _ = strconv.Atoi(row.Cells[f["EFF_ETAB"]].Value)
				apdemande.DateStatut, _ = excelToTime(row.Cells[f["DATE_STATUT"]].Value)
				// apdemande.TxPC, _ = strconv.ParseFloat(row.Cells[f["TX_PC"]].Value, 64)
				// apdemande.TxPCUnedicDares, _ = strconv.ParseFloat(row.Cells[f["TX_PC_UNEDIC_DARES"]].Value, 64)
				// apdemande.TxPCEtatDares, _ = strconv.ParseFloat(row.Cells[f["TX_PC_ETAT_DARES"]].Value, 64)
				periodStart, _ := excelToTime(row.Cells[f["DATE_DEB"]].Value)
				periodEnd, _ := excelToTime(row.Cells[f["DATE_FIN"]].Value)
				apdemande.Periode = Periode{
					Start: periodStart,
					End:   periodEnd,
				}
				apdemande.HTA, _ = strconv.ParseFloat(row.Cells[f["HTA"]].Value, 64)
				// apdemande.MTA, _ = strconv.ParseFloat(row.Cells[f["MTA"]].Value, 64)
				apdemande.EffectifAutorise, _ = strconv.Atoi(row.Cells[f["EFF_AUTO"]].Value)
				// apdemande.ProdHTAEffectif, _ = strconv.ParseFloat(row.Cells[f["PROD_HTA_EFF"]].Value, 64)
				apdemande.MotifRecoursSE, _ = strconv.Atoi(row.Cells[f["MOTIF_RECOURS_SE"]].Value)
				//apdemande.Perimetre, _ = strconv.Atoi(row.Cells[f["PERIMETRE_AP"]].Value)
				// apdemande.RecoursAnterieur, _ = strconv.Atoi(row.Cells[f["RECOURS_ANTERIEUR"]].Value)
				// apdemande.AvisCE, _ = strconv.Atoi(row.Cells[f["AVIS_CE"]].Value)
				apdemande.HeureConsommee, _ = strconv.ParseFloat(row.Cells[f["S_HEURE_CONSOM_TOT"]].Value, 64)
				// apdemande.MontantConsomme, _ = strconv.ParseFloat(row.Cells[f["S_MONTANT_CONSOM_TOT"]].Value, 64)
				apdemande.EffectifConsomme, _ = strconv.Atoi(row.Cells[f["S_EFF_CONSOM_TOT"]].Value)

				outputChannel <- &apdemande
			}

		}

		close(outputChannel)
	}()

	return outputChannel
}

func importAPDemande(c *gin.Context) {
	insertWorker, _ := c.Keys["insertEtablissement"].(chan ValueEtablissement)
	batch := c.Params.ByName("batch")
	allFiles, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	files := allFiles["apdemande"]

	for _, file := range files {
		for apdemande := range parseAPDemande(file) {
			hash := fmt.Sprintf("%x", structhash.Md5(apdemande, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: apdemande.Siret,
					Batch: map[string]Batch{
						batch: Batch{
							Compact: map[string]bool{
								"status": false,
							},
							APDemande: map[string]*APDemande{
								hash: apdemande,
							}}}}}
			insertWorker <- value
		}
	}
}

func parseAPConso(path string) chan *APConso {
	outputChannel := make(chan *APConso)

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("Error", err)
	}
	go func() {
		for _, sheet := range xlFile.Sheets {
			fields := sheet.Rows[0]
			idxID := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ID_DA" })
			idxSiret := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ETAB_SIRET" })
			idxPeriode := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MOIS" })
			idxHeureConsommee := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "HEURES" })
			idxMontants := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MONTANTS" })
			idxEffectifs := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "EFFECTIFS" })

			for _, row := range sheet.Rows[1:] {
				if len(row.Cells) > 0 {
					apconso := APConso{}
					apconso.ID = row.Cells[idxID].Value
					apconso.Siret = row.Cells[idxSiret].Value
					apconso.Periode, err = excelToTime(row.Cells[idxPeriode].Value)
					apconso.HeureConsommee, err = strconv.ParseFloat(row.Cells[idxHeureConsommee].Value, 64)
					apconso.Montant, err = strconv.ParseFloat(row.Cells[idxMontants].Value, 64)
					apconso.Effectif, err = strconv.Atoi(row.Cells[idxEffectifs].Value)

					if err != nil {
						fmt.Println(err)
					}
					outputChannel <- &apconso
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importAPConso(c *gin.Context) {
	insertWorker, _ := c.Keys["insertEtablissement"].(chan ValueEtablissement)
	batch := c.Params.ByName("batch")
	allFiles, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	files := allFiles["apconso"]

	for _, file := range files {
		for apconso := range parseAPConso(file) {
			hash := fmt.Sprintf("%x", structhash.Md5(apconso, 1))
			value := ValueEtablissement{
				Value: Etablissement{
					Siret: apconso.Siret,
					Batch: map[string]Batch{
						batch: Batch{
							APConso: map[string]*APConso{
								hash: apconso,
							}}}}}
			insertWorker <- value
		}
	}
}
