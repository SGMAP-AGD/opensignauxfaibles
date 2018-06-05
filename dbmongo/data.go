package main

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func dataDebit(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	var data interface{}

	db.C("Etablissement").Find(bson.M{"value.siret": c.Params.ByName("siret")}).One(&data)
	c.JSON(200, data)
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	// Détermination scope traitement
	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	siret := c.Params.ByName("siret")
	if siret == "" {
		queryEtablissement = bson.M{"value.index.algo1": true}
		queryEntreprise = nil
		output = bson.M{"replace": "algo1"}
	} else {
		queryEtablissement = bson.M{"value.siret": siret,
			"value.index.algo1": true}
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	mapFctEtablissement, _ := ioutil.ReadFile("js/algo1EtablissementMap.js")
	reduceFctEtablissement, _ := ioutil.ReadFile("js/algo1EtablissementReduce.js")
	finalizeFctEtablissement, _ := ioutil.ReadFile("js/algo1EtablissementFinalize.js")

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-05-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-01-01")

	scope := bson.M{"date_debut": dateDebut,
		"date_fin":          dateFin,
		"date_fin_effectif": dateFinEffectif}

	jobEtablissement := &mgo.MapReduce{
		Map:      string(mapFctEtablissement),
		Reduce:   string(reduceFctEtablissement),
		Finalize: string(finalizeFctEtablissement),
		Out:      bson.M{"replace": "MRWorkspace"},
		Scope:    scope,
	}

	_, err := db.C("Etablissement").Find(queryEtablissement).MapReduce(jobEtablissement, nil)
	if err != nil {
		c.JSON(500, err)
		return
	}

	mapFctEntreprise, _ := ioutil.ReadFile("js/algo1EntrepriseMap.js")
	reduceFctEntreprise, _ := ioutil.ReadFile("js/algo1EntrepriseReduce.js")
	finalizeFctEntreprise, _ := ioutil.ReadFile("js/algo1EntrepriseFinalize.js")

	jobEntreprise := &mgo.MapReduce{
		Map:      string(mapFctEntreprise),
		Reduce:   string(reduceFctEntreprise),
		Finalize: string(finalizeFctEntreprise),
		Out:      bson.M{"merge": "MRWorkspace"},
		Scope:    scope,
	}

	_, err = db.C("Entreprise").Find(queryEntreprise).MapReduce(jobEntreprise, nil)
	if err != nil {
		c.JSON(500, err)
	}

	mapFctUnion, _ := ioutil.ReadFile("js/algo1UnionMap.js")
	reduceFctUnion, _ := ioutil.ReadFile("js/algo1UnionReduce.js")
	finalizeFctUnion, _ := ioutil.ReadFile("js/algo1UnionFinalize.js")

	jobUnion := &mgo.MapReduce{
		Map:      string(mapFctUnion),
		Reduce:   string(reduceFctUnion),
		Finalize: string(finalizeFctUnion),
		Out:      output,
		Scope:    scope,
	}

	if output == nil {
		_, err = db.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, &result)
	} else {
		_, err = db.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, nil)
	}

	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, result)
	}
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
