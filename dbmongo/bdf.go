package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cnf/structhash"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// BDF Information Banque de France
type BDF struct {
	Siren               string    `json:"siren" bson:"siren"`
	Annee               int       `json:"annee" bson:"annee"`
	ArreteBilan         time.Time `json:"arrete_bilan" bson:"arrete_bilan"`
	RaisonSociale       string    `json:"raison_sociale" bson:"raison_sociale"`
	Secteur             string    `json:"secteur" bson:"secteur"`
	PoidsFrng           float64   `json:"poids_frng" bson:"poids_frng"`
	TauxMarge           float64   `json:"taux_marge" bson:"taux_marge"`
	DelaiFournisseur    float64   `json:"delai_fournisseur" bson:"delai_fournisseur"`
	DetteFiscale        float64   `json:"dette_fiscale" bson:"dette_fiscale"`
	FinancierCourtTerme float64   `json:"financier_court_terme" bson:"financier_court_terme"`
	FraisFinancier      float64   `json:"frais_financier" bson:"frais_financier"`
}

func parseBDF(path []string) chan *BDF {
	outputChannel := make(chan *BDF)

	go func() {
		for _, file := range path {
			xlFile, err := xlsx.OpenFile(file)

			if err != nil {
				fmt.Println("Error", err)
			}
			for _, sheet := range xlFile.Sheets {
				for _, row := range sheet.Rows[1:] {
					bdf := BDF{}
					bdf.Siren = strings.Replace(row.Cells[0].Value, " ", "", -1)
					bdf.Annee, err = strconv.Atoi(row.Cells[1].Value)
					bdf.ArreteBilan, err = excelToTime(row.Cells[2].Value)
					bdf.RaisonSociale = row.Cells[3].Value
					bdf.Secteur = row.Cells[4].Value
					bdf.PoidsFrng, err = strconv.ParseFloat(row.Cells[5].Value, 64)
					bdf.TauxMarge, err = strconv.ParseFloat(row.Cells[6].Value, 64)
					bdf.DelaiFournisseur, err = strconv.ParseFloat(row.Cells[7].Value, 64)
					bdf.DetteFiscale, err = strconv.ParseFloat(row.Cells[8].Value, 64)
					bdf.FinancierCourtTerme, err = strconv.ParseFloat(row.Cells[9].Value, 64)
					bdf.FraisFinancier, err = strconv.ParseFloat(row.Cells[10].Value, 64)

					outputChannel <- &bdf
				}
			}
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importBDF(c *gin.Context) {
	insertWorker, _ := c.Keys["insertEntreprise"].(chan ValueEntreprise)
	batch := c.Params.ByName("batch")
	allFiles, _ := GetFileList(viper.GetString("APP_DATA"), batch)
	files := allFiles["bdf"]
	for bdf := range parseBDF(files) {
		hash := fmt.Sprintf("%x", structhash.Md5(bdf, 1))

		value := ValueEntreprise{
			Value: Entreprise{
				Siren: bdf.Siren,
				Batch: map[string]Batch{
					batch: Batch{
						BDF: map[string]*BDF{
							hash: bdf,
						}}}}}
		insertWorker <- value
	}
}
