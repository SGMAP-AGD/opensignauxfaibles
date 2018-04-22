package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// Batch Lot d'intégration d'établissement
type Batch struct {
	ID        bson.ObjectId     `json:"id" bson:"_id"`
	Label     string            `json:"label" bson:"label"`
	Session   time.Time         `json:"periode" bson:"periode"`
	Ressource map[string]string `json:"ressource" bson:"ressource"`
}

func createRepo(c *gin.Context) {
	// db := c.Keys["DB"].(*mgo.Database)
	basePath := viper.GetString("APP_DATA")

	// db.C("region").Find(bson.M{"ID"})
	directories := []string{
		"admin_urssaf",
		"apconso",
		"apdemande",
		"bdf",
		"ccsf",
		"cotisation",
		"debit",
		"delais",
		"effectif",
	}

	var response map[string]string
	var status int
	for _, directory := range directories {
		path := basePath + "/" + directory
		err := os.MkdirAll(path, 700)
		status = 200
		if err != nil {
			status = 207
			response[path] = err.Error()
		} else {
			response[path] = "ok"
		}
	}
	c.JSON(status, "ok")
}

func importAP(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	go func() {
		for activitePartielleDemande := range parseActivitePartielleDemande("data-raw/activite_partielle/act_partielle_ddes_2012_janv2018.xlsx") {
			insertValue(db, activitePartielleDemande)
		}
	}()

	go func() {
		for activitePartielleConsommation := range parseActivitePartielleConsommation("data-raw/activite_partielle/act_partielle_conso_janv2018.xlsx") {
			insertValue(db, activitePartielleConsommation)
		}
	}()

	c.JSON(200, "Import done.")
}

func purge(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("Etablissement").RemoveAll(nil)
	db.C("testcollection").RemoveAll(nil)
	db.C("Debit").RemoveAll(nil)
	db.C("Delais").RemoveAll(nil)
	db.C("Cotisation").RemoveAll(nil)
	c.String(200, "Done")
}

func importCotisation(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201801.csv")

	cotisations := []string{
		"data-raw/cotisations/cotisations_bourgogne_0_201608.csv",
		"data-raw/cotisations/cotisations_bourgogne_201609_201701.csv",
		"data-raw/cotisations/cotisations_bourgogne_201702_201703.csv",
		"data-raw/cotisations/cotisations_bourgogne_201704_201707.csv",
		"data-raw/cotisations/cotisations_bourgogne_201708_201710.csv",
		"data-raw/cotisations/cotisations_bourgogne_201711.csv",
		"data-raw/cotisations/cotisations_bourgogne_201712.csv",
		"data-raw/cotisations/cotisations_bourgogne_201801.csv",
		"data-raw/cotisations/cotisations_bourgogne_201802.csv",
		"data-raw/cotisations/cotisations_frc_0_201701.csv",
		"data-raw/cotisations/cotisations_frc_201702_201703.csv",
		"data-raw/cotisations/cotisations_frc_201704_201707.csv",
		"data-raw/cotisations/cotisations_frc_201708_201710.csv",
		"data-raw/cotisations/cotisations_frc_201711.csv",
		"data-raw/cotisations/cotisations_frc_201712.csv",
		"data-raw/cotisations/cotisations_frc_201801.csv",
		"data-raw/cotisations/cotisations_frc_201802.csv",
	}

	ancienSiret := ""
	D := Value{}

	for i := range cotisations {
		for cotisation := range parseCotisation(cotisations[i], Mapping) {
			if cotisation.Value.Siret != ancienSiret {
				insertValue(db, D)
				D = Value{
					Value: Etablissement{
						Siret: cotisation.Value.Siret,
						Compte: Compte{
							Cotisation: map[string]Cotisation{},
						},
					},
				}
			}
			ancienSiret = cotisation.Value.Siret
			for k, v := range cotisation.Value.Compte.Cotisation {
				D.Value.Compte.Cotisation[k] = v
			}
		}
		insertValue(db, D)
	}
}

func getAdminUrssaf(path string) {}

// GetFileList construit la liste des fichiers à traiter
func GetFileList(basePath string, region string, period string) (map[string][]os.FileInfo, map[string]error) {
	list := make(map[string][]os.FileInfo)
	err := make(map[string]error)
	directories := []string{"admin_urssaf", "altares", "altares", "apdemande", "bdf", "ccsf", "cotisation", "debit", "delais", "effectif"}

	for _, dir := range directories {
		list[dir], err[dir] = ioutil.ReadDir(fmt.Sprintf("%s/%s/%s/%s", basePath, region, period, dir))
	}

	return list, err
}
func importAltares(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		insertValue(db, altares)
	}
}

func importEffectif(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201801.csv") {
		insertValue(db, effectif)
	}
}

func importData(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	go func() {
		for activitePartielleDemande := range parseActivitePartielleDemande("data-raw/activite_partielle/act_partielle_ddes_2012_janv2018.xlsx") {
			insertValue(db, activitePartielleDemande)
		}
	}()

	go func() {
		for activitePartielleConsommation := range parseActivitePartielleConsommation("data-raw/activite_partielle/act_partielle_conso_janv2018.xlsx") {
			insertValue(db, activitePartielleConsommation)
		}
	}()

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201801.csv") {
		insertValue(db, effectif)
	}

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		insertValue(db, altares)
	}

	go importDebit(c)
	go importCotisation(c)
	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201801.csv")

	ccsfs := []string{
		"data-raw/ccsv/Bourgogne_ccsf.csv",
		"data-raw/ccsv/FRC_ccsf.csv",
	}

	for i := range ccsfs {
		for ccsf := range parseCCSF(ccsfs[i], Mapping) {
			insertValue(db, ccsf)
		}
	}
	delaiss := []string{
		"data-raw/delais/delais_bourgogne_201301_201701.csv",
		"data-raw/delais/delais_bourgogne_201702_201707.csv",
		"data-raw/delais/delais_bourgogne_201708_201710.csv",
		"data-raw/delais/delais_bourgogne_201711.csv",
		"data-raw/delais/delais_bourgogne_201712.csv",
		"data-raw/delais/delais_bourgogne_201801.csv",
		"data-raw/delais/delais_frc_201301_201701.csv",
		"data-raw/delais/delais_frc_201702_201707.csv",
		"data-raw/delais/delais_frc_201708_201710.csv",
		"data-raw/delais/delais_frc_201711.csv",
		"data-raw/delais/delais_frc_201712.csv",
		"data-raw/delais/delais_frc_201801.csv",
	}
	for i := range delaiss {
		for delais := range parseDelais(delaiss[i], Mapping) {
			insertValue(db, delais)
		}
	}

}

func importDebit(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201801.csv")

	debits := []string{
		"data-raw/debits/debits_bourgogne_0_201607.csv",
		"data-raw/debits/debits_bourgogne_201608_201611.csv",
		"data-raw/debits/debits_bourgogne_201612_201701.csv",
		"data-raw/debits/debits_bourgogne_201702_201703.csv",
		"data-raw/debits/debits_bourgogne_201704.csv",
		"data-raw/debits/debits_bourgogne_201705_201707.csv",
		"data-raw/debits/debits_bourgogne_201708_201710.csv",
		"data-raw/debits/debits_bourgogne_201711.csv",
		"data-raw/debits/debits_bourgogne_201712.csv",
		"data-raw/debits/debits_bourgogne_201801.csv",
		"data-raw/debits/debits_bourgogne_201802.csv",
		"data-raw/debits/debits_frc_0_201611.csv",
		"data-raw/debits/debits_frc_201612_201701.csv",
		"data-raw/debits/debits_frc_201701_201703.csv",
		"data-raw/debits/debits_frc_201704.csv",
		"data-raw/debits/debits_frc_201705_201707.csv",
		"data-raw/debits/debits_frc_201708_201710.csv",
		"data-raw/debits/debits_frc_201711.csv",
		"data-raw/debits/debits_frc_201712.csv",
		"data-raw/debits/debits_frc_201801.csv",
		"data-raw/debits/debits_bourgogne_201802.csv",
	}

	ancienSiret := ""
	D := Value{}

	for i := range debits {
		for debit := range parseDebit(debits[i], Mapping) {
			if debit.Value.Siret != ancienSiret {
				insertValue(db, D)
				D = Value{
					Value: Etablissement{
						Siret: debit.Value.Siret,
						Compte: Compte{
							Debit: map[string]Debit{},
						},
					},
				}
			}
			ancienSiret = debit.Value.Siret
			for k, v := range debit.Value.Compte.Debit {
				D.Value.Compte.Debit[k] = v
			}
		}
		insertValue(db, D)
	}
}
