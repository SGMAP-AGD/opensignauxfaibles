package main

import (
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (value1 ValueEtablissement) merge(value2 ValueEtablissement) (ValueEtablissement, error) {
	if value1.Value.Siret != value2.Value.Siret {
		return ValueEtablissement{},
			errors.New("Valeurs non missibles: sirets '" +
				value1.Value.Siret + "' et '" +
				value2.Value.Siret + "'")
	}
	for idBatch := range value2.Value.Batch {
		if value1.Value.Batch == nil {
			value1.Value.Batch = make(map[string]Batch)
		}
		value1.Value.Batch[idBatch].merge(value2.Value.Batch[idBatch])
	}
	return value1, nil
}

func dataDebit(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var data interface{}

	// mapFct, _ := ioutil.ReadFile("js/dataDebit_map.js")
	// reduceFct, _ := ioutil.ReadFile("js/dataDebit_reduce.js")
	// finalizeFct, _ := ioutil.ReadFile("js/dataDebit_finalize.js")

	// job := &mgo.MapReduce{
	// 	Map:      string(mapFct),
	// 	Reduce:   string(reduceFct),
	// 	Finalize: string(finalizeFct),
	// }

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).One(&data)
	c.JSON(200, data)
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/algo1_map.js")
	reduceFct, _ := ioutil.ReadFile("js/algo1_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/algo1_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
	}

	var etablissement []interface{}

	db.C("Entreprise").Find(bson.M{"value.siren": c.Params.ByName("siren")}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func reduceAll(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	mapFct, _ := ioutil.ReadFile("js/algo1_map.js")
	reduceFct, _ := ioutil.ReadFile("js/algo1_reduce.js")
	finalizeFct, _ := ioutil.ReadFile("js/algo1_finalize.js")

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      bson.M{"replace": "algo1"},
	}

	var etablissement []struct {
		ID    string      `json:"id" bson:"_id"`
		Value interface{} `json:"value" bson:"value"`
	}

	db.C("Entreprise").Find(bson.M{"value.index.algo1": true}).MapReduce(job, &etablissement)

	c.JSON(200, etablissement)
}

func browse(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var entreprise []interface{}
	db.C("Entreprise").Find(bson.M{"value.siren": c.Params.ByName("siren")}).All(&entreprise)
	c.JSON(200, entreprise)
}

func browseOrig(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	var etablissement []interface{}
	db.C("Etablissement").Find(bson.M{"siret": c.Params.ByName("siret")}).All(&etablissement)
	c.JSON(200, etablissement)
}

func compactEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}

	siret := c.Params.ByName("siret")
	if siret == "" {
		query = nil
		output = bson.M{"replace": "Etablissement"}
		etablissement = nil
	} else {
		query = bson.M{"value.siret": siret}
		output = nil
	}

	// Ressources JS
	mapFct, errMap := ioutil.ReadFile("js/compactEtablissementMap.js")
	reduceFct, errReduce := ioutil.ReadFile("js/compactEtablissementReduce.js")
	finalizeFct, errFinalize := ioutil.ReadFile("js/compactEtablissementFinalize.js")
	if errMap != nil || errReduce != nil || errFinalize != nil {
		c.JSON(500, "Impossible d'accéder aux ressources JS pour ce traitement: "+errMap.Error()+" "+errFinalize.Error()+" "+errReduce.Error())
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      output,
		Scope: bson.M{"batches": []string{"1802", "1803", "1804", "1805"},
			"types": []string{
				"altares",
				"apconso",
				"apdemande",
				"ccsf",
				"cotisation",
				"debit",
				"delai",
				"effectif",
				"sirene",
				"dpae",
			},
			"deleteOld": []string{"effectif", "apdemande", "apconso", "altares"},
		},
	}

	err := errors.New("")
	if output == nil {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, "Echec du traitement MR, message serveur: "+err.Error())
	} else {
		c.JSON(200, etablissement)
	}

}

func compactEntreprise(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}

	siren := c.Params.ByName("siren")
	if siren == "" {
		query = nil
		output = bson.M{"replace": "Entreprise"}
		etablissement = nil
	} else {
		query = bson.M{"value.siren": siren}
		output = nil
	}

	// Ressources JS
	mapFct, errMap := ioutil.ReadFile("js/compactEntrepriseMap.js")
	reduceFct, errReduce := ioutil.ReadFile("js/compactEntrepriseReduce.js")
	finalizeFct, errFinalize := ioutil.ReadFile("js/compactEntrepriseFinalize.js")
	if errMap != nil || errReduce != nil || errFinalize != nil {
		c.JSON(500, "Impossible d'accéder aux ressources JS pour ce traitement: "+errMap.Error()+" "+errFinalize.Error()+" "+errReduce.Error())
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      output,
		Scope: bson.M{"batches": []string{"1802", "1803", "1804", "1805"},
			"types": []string{
				"bdf",
			},
			"deleteOld": []string{"bdf"},
		},
	}

	err := errors.New("")
	if output == nil {
		_, err = db.C("Entreprise").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.C("Entreprise").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, "Echec du traitement MR, message serveur: "+err.Error())
	} else {
		c.JSON(200, etablissement)
	}

}
