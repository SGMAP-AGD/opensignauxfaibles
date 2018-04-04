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

	declareServerFunctions(db)
	if err != nil {
		log.Panic(err)
	}

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func insertValue(db *mgo.Database, value Value) {
	value.ID = bson.NewObjectId()
	db.C("Etablissement").Insert(value)
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

func declareServerFunctions(db *mgo.Database) {
	f := ServerJSFunc{
		ID:    "generatePeriodSerie",
		Value: "function (date_debut, date_fin) {var date_next = new Date(date_debut.getTime());var serie = [];while (date_next.getTime() < date_fin.getTime()) {serie.push(new Date(date_next.getTime()));date_next.setUTCMonth(date_next.getUTCMonth() + 1);}return serie;",
	}
	f.Add(db)
	f = ServerJSFunc{
		ID:    "compareDebit",
		Value: `function(a,b) {if (a.numero_historique < b.numero_historique) return -1;if (a.numero_historique > b.numero_historique) return 1;return 0;}`,
	}
	f.Add(db)
	f = ServerJSFunc{
		ID:    "isRJLJ",
		Value: `function(code) {codes = ["PCL010501","PCL010502","PCL030105","PCL05010102","PCL05010203","PCL05010402","PCL05010302","PCL05010502","PCL05010702","PCL05010802","PCL05010901","PCL05011003","PCL05011101","PCL05011203","PCL05011303","PCL05011403","PCL05011503","PCL05011603","PCL05011902","PCL05012003","PCL0108","PCL0109","PCL030107","PCL030108","PCL030307","PCL030308","PCL05010103","PCL05010104","PCL05010204","PCL05010205","PCL05010303","PCL05010304","PCL05010403","PCL05010404","PCL05010503","PCL05010504","PCL05010703","PCL05010803","PCL05011004","PCL05011005","PCL05011102","PCL05011103","PCL05011204","PCL05011205","PCL05011304","PCL05011305","PCL05011404","PCL05011405","PCL05011504","PCL05011505","PCL05011604","PCL05011605","PCL05011903","PCL05011904","PCL05012004","PCL05012005","PCL040802"];return codes.includes(code);}`,
	}
	f.Add(db)
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

	db.C("Etablissement").Find(bson.M{"_id": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

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
		Out:      bson.M{"replace": "algo1"},
	}

	var etablissement []struct {
		ID    Siret       `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("Etablissement").Find(bson.M{"value.index.algo1": true}).MapReduce(job, &etablissement)

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

func importAltares(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		insertValue(db, altares)
	}
}

func importEffectif(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv") {
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

	for effectif := range parseEffectif("data-raw/effectif/Urssaf_emploi_BFC_201001_201709.csv") {
		insertValue(db, effectif)
	}

	for altares := range parseAltares("data-raw/altares/RECAP_ALTARES_201803.csv") {
		insertValue(db, altares)
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

	mapFct, _ := ioutil.ReadFile("map.js")
	reduceFct, _ := ioutil.ReadFile("reduce.js")
	finalizeFct, _ := ioutil.ReadFile("finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var etablissement []interface{}

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceEtablissements(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("map.js")
	reduceFct, _ := ioutil.ReadFile("reduce.js")
	finalizeFct, _ := ioutil.ReadFile("finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "Etablissement"},
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
