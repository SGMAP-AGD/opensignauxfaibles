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

// Sirene informations sur les entreprises
type Sirene struct {
	Siren              string    `json,omitempty:"siren" bson,omitempty:"siren"`
	Nic                string    `json,omitempty:"nic" bson,omitempty:"nic"`
	RaisonSociale      string    `json,omitempty:"raison_sociale" bson,omitempty:"raison_sociale"`
	NumVoie            string    `json,omitempty:"numero_voie" bson,omitempty:"numero_voie"`
	IndRep             string    `json,omitempty:"indrep" bson,omitempty:"indrep"`
	TypeVoie           string    `json,omitempty:"type_voie" bson,omitempty:"type_voie"`
	CodePostal         string    `json,omitempty:"code_postal" bson,omitempty:"code_postal"`
	Cedex              string    `json,omitempty:"cedex" bson,omitempty:"cedex"`
	Region             string    `json,omitempty:"region" bson,omitempty:"region"`
	Departement        string    `json,omitempty:"departement" bson,omitempty:"departement"`
	Commune            string    `json,omitempty:"commune" bson,omitempty:"commune"`
	APE                string    `json,omitempty:"ape" bson,omitempty:"ape"`
	NatureActivite     string    `json,omitempty:"" bson,omitempty:""`
	ActiviteSaisoniere string    `json,omitempty:"" bson,omitempty:""`
	ModaliteActivite   string    `json,omitempty:"modalite_activite" bson,omitempty:"modalite_activite"`
	Productif          string    `json,omitempty:"productif" bson,omitempty:"productif"`
	NatureJuridique    string    `json,omitempty:"nature_juridique" bson,omitempty:"nature_juridique"`
	Categorie          string    `json,omitempty:"categorie" bson,omitempty:"categorie"`
	Creation           time.Time `json,omitempty:"date_creation" bson,omitempty:"date_creation"`
	IndiceMonoactivite int       `json,omitempty:"indice_monoactivite" bson,omitempty:"indice_monoactivite"`
	TrancheCA          int       `json,omitempty:"tranche_ca" bson,omitempty:"tranche_ca"`
	Sigle              string    `json,omitempty:"sigle" bson,omitempty:"sigle"`
	DebutActivite      time.Time `json:"debut_activite" bson:"debut_activite"`
}

func parseSirene(paths []string, batch string) chan Sirene {
	outputChannel := make(chan Sirene)
	go func() {
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'

			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					log.Fatal(error)
				}

				sirene := Sirene{}
				sirene.Siren = row[0]
				sirene.Nic = row[1]
				sirene.RaisonSociale = row[36]
				sirene.NumVoie = row[16]
				sirene.IndRep = row[17]
				sirene.TypeVoie = row[19]
				sirene.CodePostal = row[20]
				sirene.Cedex = row[21]
				sirene.Region = row[23]
				sirene.Departement = row[24]
				sirene.Commune = row[28]
				sirene.APE = row[42]
				sirene.NatureActivite = row[52]
				sirene.ActiviteSaisoniere = row[55]
				sirene.ModaliteActivite = row[56]
				sirene.Productif = row[57]
				sirene.NatureJuridique = row[71]
				sirene.Categorie = row[82]
				sirene.Creation, _ = time.Parse("20060102", row[50])
				sirene.IndiceMonoactivite, _ = strconv.Atoi(row[85])
				sirene.TrancheCA, _ = strconv.Atoi(row[89])
				sirene.Sigle = row[61]
				sirene.DebutActivite, _ = time.Parse("20060102", row[51])

				outputChannel <- sirene
			}
			file.Close()
		}
		close(outputChannel)
	}()

	return outputChannel
}

// hash := fmt.Sprintf("%x", structhash.Md5(sirene, 1))

func importSirene(c *gin.Context) {
	insertWorker := c.Keys["DBW"].(chan Value)
	batch := c.Params.ByName("batch")

	files, _ := GetFileList(viper.GetString("APP_DATA"), batch)

	sirene := files["sirene"]
	for sirene := range parseSirene(sirene, batch) {
		hash := fmt.Sprintf("%x", structhash.Md5(sirene, 1))

		value := Value{
			Value: Entreprise{
				Siren:  sirene.Siren,
				Region: sirene.Region,
				Etablissement: map[string]Etablissement{
					sirene.Siren + sirene.Nic: Etablissement{
						Siret: sirene.Siren + sirene.Nic,
						Batch: map[string]Batch{
							batch: Batch{
								Compact: map[string]bool{
									"status": false,
								},
								Sirene: map[string]Sirene{
									hash: sirene,
								}}}}}}}
		insertWorker <- value
	}

	insertWorker <- Value{}

}
