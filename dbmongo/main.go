package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DB Initialisation de la connexion MongoDB
func DB() gin.HandlerFunc {

	mongodb, err := mgo.Dial("127.0.0.1")
	db := mongodb.DB("jason")

	if err != nil {
		log.Panic(err)
	}

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(DB())
	v1 := r.Group("api/v1")
	{
		v1.GET("/purge", purge)
		v1.GET("/import", importData)
		v1.GET("/reduceEtablissement/:siret", reduceEtablissement)
		v1.GET("/reduceEtablissement", reduceEtablissements)
		v1.GET("/reduce/:siret", reduce)
		v1.GET("/reduce", reduceAll)
		v1.GET("/etablissement/:siret", browseEtablissement)
		v1.GET("/orig/:siret", browseOrig)
		v1.GET("/debug/:urssaf", debug)
		v1.GET("/importAP", importAP)
		v1.GET("/importDebit", importDebit)
		v1.GET("/importAltares", importAltares)
		v1.GET("/mapDebit", mapDebit)
		v1.GET("/importEffectif", importEffectif)
	}

	r.Run(":8080")
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map2.js")
	reduceFct, _ := ioutil.ReadFile("reduce2.js")
	finalizeFct, _ := ioutil.ReadFile("finalize2.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		//Out:    bson.M{"replace": "testcollection2"},
	}

	var etablissement []struct {
		ID    Siret       `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("testcollection").Find(bson.M{"_id": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceAll(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map2.js")
	reduceFct, _ := ioutil.ReadFile("reduce2.js")
	finalizeFct, _ := ioutil.ReadFile("finalize2.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "testcollection2"},
	}

	var etablissement []struct {
		ID    Siret       `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("testcollection").Find(nil).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func mapDebit(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map_debit.js")
	reduceFct, _ := ioutil.ReadFile("reduce2.js")

	job := &mgo.MapReduce{
		Map:    string(mapFct),
		Reduce: string(reduceFct),
		//Out:    bson.M{"replace": "testcollection2"},
	}

	var etablissement []struct {
		ID    Siret       `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("testcollection").Find(bson.M{"_id": "80969365800027"}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func debug(c *gin.Context) {
	ursaff := c.Params.ByName("urssaf")
	date, _ := UrssafToPeriod(ursaff)
	c.JSON(200, date)
}

func importDebit(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv")

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
	D := Etablissement{}

	for i := range debits {
		for debit := range parseDebit(debits[i], Mapping) {
			if debit.Siret != ancienSiret {
				db.C("Etablissement").Insert(D)
				D = Etablissement{
					Compte: Compte{
						Debit: map[string]Debit{},
					},
				}
				D.Siret = debit.Siret
			}
			ancienSiret = debit.Siret
			for k, v := range debit.Compte.Debit {
				D.Compte.Debit[k] = v
			}
		}
		db.C("Etablissement").Insert(D)
	}
}

func importAltares(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		db.C("Etablissement").Insert(altares)
	}
}

func importEffectif(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv") {
		db.C("Etablissement").Insert(effectif)
	}
}

func importData(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	go func() {
		for activitePartielleDemande := range parseActivitePartielleDemande("data-raw/activite_partielle/act_partielle_ddes_2012_janv2018.xlsx") {
			db.C("Etablissement").Insert(activitePartielleDemande)
		}
	}()

	go func() {
		for activitePartielleConsommation := range parseActivitePartielleConsommation("data-raw/activite_partielle/act_partielle_conso_janv2018.xlsx") {
			db.C("Etablissement").Insert(activitePartielleConsommation)
		}
	}()

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv") {
		db.C("Etablissement").Insert(effectif)
	}

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		db.C("Etablissement").Insert(altares)
	}

	go importDebit(c)
	go importCotisation(c)
	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv")

	ccsfs := []string{
		"data-raw/ccsv/Bourgogne_ccsf.csv",
		"data-raw/ccsv/FRC_ccsf.csv",
	}

	for i := range ccsfs {
		for ccsf := range parseCCSF(ccsfs[i], Mapping) {
			db.C("Etablissement").Insert(ccsf)
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
			db.C("Etablissement").Insert(delais)
		}
	}

}

func importCotisation(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	Mapping := getCompteSiretMapping("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv")

	cotisations := []string{
		"data-raw/cotisations/cotisations_bourgogne_0_201608.csv",
		"data-raw/cotisations/cotisations_bourgogne_201609_201701.csv",
		"data-raw/cotisations/cotisations_bourgogne_201702_201703.csv",
		"data-raw/cotisations/cotisations_bourgogne_201704_201707.csv",
		"data-raw/cotisations/cotisations_bourgogne_201708_201710.csv",
		"data-raw/cotisations/cotisations_bourgogne_201711.csv",
		"data-raw/cotisations/cotisations_bourgogne_201712.csv",
		"data-raw/cotisations/cotisations_bourgogne_201801.csv",
		"data-raw/cotisations/cotisations_frc_0_201701.csv",
		"data-raw/cotisations/cotisations_frc_201702_201703.csv",
		"data-raw/cotisations/cotisations_frc_201704_201707.csv",
		"data-raw/cotisations/cotisations_frc_201708_201710.csv",
		"data-raw/cotisations/cotisations_frc_201711.csv",
		"data-raw/cotisations/cotisations_frc_201712.csv",
		"data-raw/cotisations/cotisations_frc_201801.csv",
	}

	ancienSiret := ""
	D := Etablissement{}

	for i := range cotisations {
		for cotisation := range parseCotisation(cotisations[i], Mapping) {
			if cotisation.Siret != ancienSiret {
				db.C("Etablissement").Insert(D)
				D = Etablissement{
					Compte: Compte{
						Cotisation: map[string]Cotisation{},
					},
				}
				D.Siret = cotisation.Siret
			}
			ancienSiret = cotisation.Siret
			for k, v := range cotisation.Compte.Cotisation {
				D.Compte.Cotisation[k] = v
			}
		}
		db.C("Etablissement").Insert(D)
	}
}

func importAP(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	go func() {
		for activitePartielleDemande := range parseActivitePartielleDemande("data-raw/activite_partielle/act_partielle_ddes_2012_janv2018.xlsx") {
			db.C("Etablissement").Insert(activitePartielleDemande)
		}
	}()

	go func() {
		for activitePartielleConsommation := range parseActivitePartielleConsommation("data-raw/activite_partielle/act_partielle_conso_janv2018.xlsx") {
			db.C("Etablissement").Insert(activitePartielleConsommation)
		}
	}()

	c.JSON(200, "Import done.")
}

func browseEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var etablissement []interface{}
	db.C("testcollection").Find(bson.M{"_id": c.Params.ByName("siret")}).All(&etablissement)
	c.JSON(200, etablissement)
}

func browseOrig(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var etablissement []interface{}
	db.C("Etablissement").Find(bson.M{"siret": c.Params.ByName("siret")}).All(&etablissement)
	c.JSON(200, etablissement)
}

func reduceEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	reduceFct, _ := ioutil.ReadFile("reduce.js")
	finalizeFct, _ := ioutil.ReadFile("finalize.js")

	job := &mgo.MapReduce{
		Map:      "function() { emit(this.siret, this) }",
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		// Out:    bson.M{"replace": "testcollection"},
	}

	var etablissement []struct {
		ID    Siret         `json:"id" bson:"_id"`
		Value Etablissement `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(bson.M{"siret": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceEtablissements(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	reduceFct, _ := ioutil.ReadFile("reduce.js")

	job := &mgo.MapReduce{
		Map:    "function() { emit(this.siret, this) }",
		Reduce: string(reduceFct),
		Out:    bson.M{"replace": "testcollection"},
	}

	var etablissement []struct {
		ID    Siret         `json:"id" bson:"_id"`
		Value Etablissement `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(nil).MapReduce(job, nil)

	c.JSON(200, etablissement)
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
