package main

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func data(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	siret := c.Params.ByName("siret")
	var query interface{}
	var output interface{}

	if siret == "" {
		query = nil
		output = bson.M{"replace": data}
	} else {
		query = bson.M{"value.siret": bson.M{"$in": []string{"31039150300029",
			"32280493100044",
			"32411592200043",
			"32412727300013",
			"32905885300022",
			"33050826800019",
			"33468704300052",
			"33973099600024",
			"34285177100010",
			"34763185500017",
			"34853799400017",
			"35137354300039",
			"37908891700016",
			"37944138900028",
			"38062893300034",
			"38165327800089",
			"38389828500010",
			"38900298100030",
			"39020638100019",
			"39359780200064",
			"39361400300050",
			"39385363500026",
			"39749044200050",
			"39829078300024",
			"40069600100058",
			"40197001700026",
			"40843381100036",
			"40848145500017",
			"41034282800036",
			"41091104400056",
			"41221709300019",
			"41484267400015",
			"41809276300030",
			"41827199500056",
			"41873205300024",
			"41902272800010",
			"42054736600047",
			"42072807300016",
			"42269962900024",
			"43008570400012",
			"44030310500025",
			"44829364700013",
			"45136355000018",
			"45277190000027",
			"48322898700010",
			"49310588600011",
			"49501449000017",
			"50695021100025",
			"62558027900127",
			"65715007400026",
			"67725020100022",
			"68282031100038",
			"70558010800011",
			"72262116600049",
			"77829308400068",
			"77857784100019",
			"79712014400010"}}}
		output = nil
	}

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-05-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-01-01")

	scope := bson.M{"date_debut": dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
	}

	mapFct, errM := ioutil.ReadFile("js/dataMap.js")
	reduceFct, errR := ioutil.ReadFile("js/dataReduce.js")
	finalizeFct, errF := ioutil.ReadFile("js/dataFinalize.js")

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

	// Détermination scope traitement
	algo := c.Params.ByName("algo")
	siret := c.Params.ByName("siret")

	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	if siret == "" {
		queryEtablissement = bson.M{"value.index." + algo: true}
		queryEntreprise = nil
		output = bson.M{"replace": algo}
	} else {
		queryEtablissement = bson.M{"value.siret": siret}
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-05-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-01-01")

	mapFctEtablissement, errEtabM := ioutil.ReadFile("js/" + algo + "EtablissementMap.js")
	reduceFctEtablissement, errEtabR := ioutil.ReadFile("js/" + algo + "EtablissementReduce.js")
	finalizeFctEtablissement, errEtabF := ioutil.ReadFile("js/" + algo + "EtablissementFinalize.js")

	mapFctEntreprise, errEntM := ioutil.ReadFile("js/" + algo + "EntrepriseMap.js")
	reduceFctEntreprise, errEntR := ioutil.ReadFile("js/" + algo + "EntrepriseReduce.js")
	finalizeFctEntreprise, errEntF := ioutil.ReadFile("js/" + algo + "EntrepriseFinalize.js")

	mapFctUnion, errUnM := ioutil.ReadFile("js/" + algo + "UnionMap.js")
	reduceFctUnion, errUnR := ioutil.ReadFile("js/" + algo + "UnionReduce.js")
	finalizeFctUnion, errUnF := ioutil.ReadFile("js/" + algo + "UnionFinalize.js")

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
