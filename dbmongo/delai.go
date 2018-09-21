package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cnf/structhash"
)

// Delai tuple fichier ursaff
type Delai struct {
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheanche     time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	NumeroStructure   string    `json:"numero_structure" bson:"numero_structure"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

func parseDelai(paths []string) chan *Delai {
	outputChannel := make(chan *Delai)

	field := map[string]int{
		"NumeroCompte":      0,
		"NumeroContentieux": 1,
		"DateCreation":      2,
		"DateEcheanche":     3,
		"DureeDelai":        4,
		"Denomination":      5,
		"Indic6m":           6,
		"AnneeCreation":     7,
		"MontantEcheancier": 8,
		"NumeroStructure":   9,
		"Stade":             10,
		"Action":            11,
	}

	go func() {
		for _, path := range paths {
			log(info, "importDelais", "Import du fichier délais "+path)
			file, err := os.Open(path)
			if err != nil {
				log(critical, "importDelais", "Erreur à l'ouverture du fichier "+path+": "+err.Error()+", abandon")
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()
			n := 0
			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log(critical, "importDelais", "Erreur à la lecture du fichier '"+path+"': «"+err.Error()+"», abandon")
					break
				}
				n++
				delai := Delai{}
				delai.NumeroCompte = row[field["NumeroCompte"]]
				delai.NumeroContentieux = row[field["NumeroContentieux"]]
				delai.DateCreation, err = time.Parse("2006-01-02", row[field["DateCreation"]])
				delai.DateEcheanche, err = time.Parse("2006-01-02", row[field["DateEcheanche"]])
				delai.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
				delai.Denomination = row[field["Denomination"]]
				delai.Indic6m = row[field["Indic6m"]]
				delai.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
				delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[field["MontantEcheancier"]], ",", ".", -1), 64)
				delai.NumeroStructure = row[field["NumeroStructure"]]
				delai.Stade = row[field["Stade"]]
				delai.Action = row[field["Action"]]

				outputChannel <- &delai
			}
			file.Close()
			log(critical, "importDelais", "Import du fichier délais "+path+" terminé"+fmt.Sprint(n)+" lignes insérées")
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importDelai(batch *AdminBatch) error {
	log(info, "importDelai", "Import du batch "+batch.ID.Key+": Délai")
	mapping, err := getCompteSiretMapping(batch)
	if err != nil {
		log(critical, "importDelai", "Erreur d'accès au mapping Siret/Compte, interruption. "+err.Error())
	}

	for delai := range parseDelai(batch.Files["delai"]) {
		if siret, ok := mapping[delai.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(delai, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							Delai: map[string]*Delai{
								hash: delai,
							}}}}}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	log(info, "importDelai", "Importation du batch "+batch.ID.Key+" terminée")
	return nil
}
