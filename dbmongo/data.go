package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func dataPrediction(c *gin.Context) {
	var prediction []Prediction
	var etablissement []ValueEtablissement
	var siret []string

	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("prediction").Find(nil).Sort("-prob").Limit(50).All(&prediction)
	for _, r := range prediction {
		siret = append(siret, r.Siret)
	}

	query := bson.M{"_id": bson.M{"$in": siret}}
	db.C("Etablissement").Find(query).All(&etablissement)
	c.JSON(200, bson.M{"prediction": prediction, "etablissement": etablissement})

}

func data(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	siret := []string{}
	c.Bind(&siret)
	fmt.Println(siret)

	var query interface{}
	var output interface{}

	query = bson.M{"value.siret": bson.M{"$in": siret}}
	output = nil

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-06-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-03-01")

	scope := bson.M{"date_debut": dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
	}

	mapFct, errM := ioutil.ReadFile("js/browse/Map.js")
	reduceFct, errR := ioutil.ReadFile("js/browse/Reduce.js")
	finalizeFct, errF := ioutil.ReadFile("js/browse/Finalize.js")

	if errM != nil || errR != nil || errF != nil {
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	job := &mgo.MapReduce{
		Map:      string(mapFct),
		Reduce:   string(reduceFct),
		Finalize: string(finalizeFct),
		Out:      nil,
		Scope:    scope,
	}

	var result interface{}
	var err error

	if output == nil {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, &result)
	} else {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, result)
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-06-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-03-01")

	// Détermination scope traitement
	algo := c.Params.ByName("algo")
	batch := c.Params.ByName("batch")
	siret := c.Params.ByName("siret")

	db.C("Features").RemoveAll(bson.M{"_id.batch": batch, "_id.algo": algo})

	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	if siret == "" {
		queryEtablissement = bson.M{"value.index." + algo: true}
		queryEntreprise = nil
		output = bson.M{"merge": "Features"}
	} else {
		queryEtablissement = bson.M{"value.siret": siret}
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	mapFctEtablissement, errEtabM := ioutil.ReadFile("js/" + algo + "/EtablissementMap.js")
	reduceFctEtablissement, errEtabR := ioutil.ReadFile("js/" + algo + "/EtablissementReduce.js")
	finalizeFctEtablissement, errEtabF := ioutil.ReadFile("js/" + algo + "/EtablissementFinalize.js")

	mapFctEntreprise, errEntM := ioutil.ReadFile("js/" + algo + "/EntrepriseMap.js")
	reduceFctEntreprise, errEntR := ioutil.ReadFile("js/" + algo + "/EntrepriseReduce.js")
	finalizeFctEntreprise, errEntF := ioutil.ReadFile("js/" + algo + "/EntrepriseFinalize.js")

	mapFctUnion, errUnM := ioutil.ReadFile("js/" + algo + "/UnionMap.js")
	reduceFctUnion, errUnR := ioutil.ReadFile("js/" + algo + "/UnionReduce.js")
	finalizeFctUnion, errUnF := ioutil.ReadFile("js/" + algo + "/UnionFinalize.js")

	if errEtabM != nil || errEtabR != nil || errEtabF != nil ||
		errEntM != nil || errEntR != nil || errEntF != nil ||
		errUnM != nil || errUnR != nil || errUnF != nil {
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	scope := bson.M{"date_debut": dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
		"actual_batch":           batch,
	}

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
		return
	}

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
		return
	}

	err = db.C("MRWorkspace").DropCollection()
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, result)

}

func browse(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	var siret []string
	c.Bind(&siret)

	var etablissement []interface{}
	db.C("Etablissement").Find(bson.M{"value.siret": "41191670300019"}).All(&etablissement)
	spew.Dump(etablissement)
	c.JSON(200, etablissement)
}

func indexEntreprise(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var entreprise []interface{}

	siret := c.Params.ByName("siren")
	if siret == "" {
		query = nil
		output = bson.M{"replace": "Etablissement"}
		entreprise = nil
	} else {
		query = bson.M{"value.siret": siret}
		output = nil
	}

	// Ressources JS
	mapFct, errMap := ioutil.ReadFile("js/index/EtablissementMap.js")
	reduceFct, errReduce := ioutil.ReadFile("js/index/EtablissementReduce.js")
	if errMap != nil || errReduce != nil {
		c.JSON(500, "Impossible d'accéder aux ressources JS pour ce traitement: "+errMap.Error()+" "+errReduce.Error())
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:    string(mapFct),
		Reduce: string(reduceFct),
		Out:    output,
		Scope: bson.M{"batches": []string{"1802", "1803", "1804", "1805", "1806"},
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
		_, err = db.C("Etablissement").Find(query).MapReduce(job, &entreprise)
	} else {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, "Echec du traitement MR, message serveur: "+err.Error())
	} else {
		c.JSON(200, entreprise)
	}
}

func compactEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batches := getBatchesID(db)

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
	mapFct, errMap := ioutil.ReadFile("js/compact/EtablissementMap.js")
	reduceFct, errReduce := ioutil.ReadFile("js/compact/EtablissementReduce.js")
	finalizeFct, errFinalize := ioutil.ReadFile("js/compact/EtablissementFinalize.js")
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
		Scope: bson.M{"batches": batches,
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
	batches := getBatchesID(db)

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
	mapFct, errMap := ioutil.ReadFile("js/compact/EntrepriseMap.js")
	reduceFct, errReduce := ioutil.ReadFile("js/compact/EntrepriseReduce.js")
	finalizeFct, errFinalize := ioutil.ReadFile("js/compact/EntrepriseFinalize.js")
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
		Scope: bson.M{"batches": batches,
			"types": []string{
				"bdf",
			},
			"deleteOld": []string{"bdf"},
		},
	}

	var err error

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
