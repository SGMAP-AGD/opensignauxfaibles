package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// prepareMRJob charge les fichiers MapReduce et fournit les paramètres pour l'exécution
func prepareMRJob(batchKey string, id string, typeJob string, target string) (*mgo.MapReduce, *bson.M, error) {
	query := &bson.M{
		"_id": id,
	}
	var output interface{}
	output = nil
	if id == "" {
		query = nil
		output = &bson.M{"merge": "Browser"}
	}

	MR := MapReduceJS{}
	errEt := MR.load(typeJob, target)
	if errEt != nil {
		return &mgo.MapReduce{}, nil, fmt.Errorf("Problème d'accès aux fichiers MapReduce " + typeJob + "/" + target)
	}

	batch, _ := getBatch(batchKey)

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	naf, errNAF := loadNAF()
	if errNAF != nil {
		return &mgo.MapReduce{}, nil, fmt.Errorf("Problème d'accès aux fichiers naf")
	}

	job := &mgo.MapReduce{
		Map:      string(MR.Map),
		Reduce:   string(MR.Reduce),
		Finalize: string(MR.Finalize),
		Out:      output,
		Scope: bson.M{
			"date_debut":             dateDebut,
			"date_fin":               dateFin,
			"date_fin_effectif":      dateFinEffectif,
			"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
			"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
			"actual_batch":           batch.ID.Key,
			"naf":                    naf,
		},
	}

	return job, query, nil
}

// browserIndexHandler traite les requêtes d'indexation
func browserIndexHandler(c *gin.Context) {
	var index struct {
		BatchKey string `json:"batch" binding:"required"`
		Siret    string `json:"siret"`
	}

	err := c.ShouldBind(&index)
	if err != nil {
		c.JSON(500, "Requête malformée: "+err.Error())
	}

	result, err := browserIndex(index.BatchKey, index.Siret)

	if err != nil {
		c.JSON(500, "Erreur du traitement: "+err.Error())
		return
	}
	c.JSON(200, result)
}

func browserIndex(batchKey string, siret string) (interface{}, error) {
	var siren string
	if len(siret) == 14 {
		siren = siret[0:9]
	}

	// préparation des jobs
	var jobs [2]*mgo.MapReduce
	var queries [2]*bson.M
	var errMR [2]error
	jobs[0], queries[0], errMR[0] = prepareMRJob(batchKey, siret, "browser", "etablissement")
	jobs[1], queries[1], errMR[1] = prepareMRJob(batchKey, siren, "browser", "entreprise")

	if !allErrors(errMR[:], nil) {
		return nil, errors.New("Erreur dans la création du job MapReduce: ")
	}

	// exécution
	var resultEtablissement interface{}
	var resultEntreprise interface{}
	var err [2]error

	_, err[0] = db.DB.C("Etablissement").Find(queries[0]).MapReduce(jobs[0], &resultEtablissement)
	_, err[1] = db.DB.C("Entreprise").Find(queries[1]).MapReduce(jobs[1], &resultEntreprise)
	if !allErrors(err[:], nil) {
		errorMessage := ""
		if err[0] != nil {
			errorMessage = errorMessage + "\nEtablissement : " + err[0].Error() + " "
		}
		if err[1] != nil {
			errorMessage = errorMessage + "\nEntreprise: " + err[1].Error()
		}
		return nil, errors.New("Erreur dans l'exécution des jobs MapReduce" + errorMessage)
	}

	return map[string]interface{}{
		"etablissement": resultEtablissement,
		"entreprise":    resultEntreprise,
	}, nil
}

func predictionBrowse(c *gin.Context) {
	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$limit": 50})
	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Entreprise",
		"localField":   "siren",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{"bdf": bson.M{"$arrayElemAt": []interface{}{"$entreprise.value.batch.1802.bdf", 0}}}})

	pipeline = append(pipeline, bson.M{"$project": bson.M{"entreprise": 0}})

	var result []interface{}
	err := db.DB.C("Prediction").Pipe(pipeline).All(&result)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, result)
}

// func predictionBrowse(c *gin.Context) {
// 	batch := c.Params.ByName("batch")
// 	algo := c.Params.ByName("algo")
// 	paramPage := c.Params.ByName("page")
// 	page, err := strconv.Atoi(paramPage)

// 	if err != nil {
// 		c.JSON(500, err)
// 		return
// 	}

// 	//siret := c.Params.ByName("siret")
// 	var pipeline []bson.M

// 	pipeline = append(pipeline, bson.M{"$match": bson.M{"_id.batch": batch, "_id.algo": algo}})
// 	pipeline = append(pipeline, bson.M{"$sort": bson.M{"score": -1}})
// 	pipeline = append(pipeline, bson.M{"$skip": 10 * page})
// 	pipeline = append(pipeline, bson.M{"$limit": 10})
// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
// 	}})

// 	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
// 		"from":         "Etablissement",
// 		"localField":   "_id.siret",
// 		"foreignField": "_id",
// 		"as":           "etablissement"}})

// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"etablissement": bson.M{"$arrayElemAt": []interface{}{"$etablissement", 0}}}})
// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"etablissement": "$etablissement.value"}})

// 	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
// 		"from":         "Entreprise",
// 		"localField":   "siren",
// 		"foreignField": "_id",
// 		"as":           "entreprise"}})

// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"entreprise": bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}}}})
// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"entreprise": "$entreprise.value"}})

// 	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
// 		"from": "Features",
// 		"let": bson.M{"siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
// 			"siret": "$_id.siret"},
// 		"pipeline": []interface{}{
// 			bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
// 				bson.M{"$eq": []interface{}{"$_id.siren", "$$siren"}},
// 				bson.M{"$eq": []interface{}{"$_id.batch", batch}},
// 				bson.M{"$eq": []interface{}{"$_id.algo", algo}},
// 			}}}},
// 			bson.M{"$addFields": bson.M{
// 				"features": bson.M{"$filter": bson.M{"input": "$value",
// 					"cond": bson.M{
// 						"$eq": []interface{}{"$$this.siret", "$$siret"}}}}}},
// 			bson.M{"$project": bson.M{
// 				"features": "$features",
// 			}},
// 		},
// 		"as": "features"}})
// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"features": bson.M{"$arrayElemAt": []interface{}{"$features", 0}}}})
// 	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
// 		"features": "$features.features"}})

// 	var result []interface{}
// 	err = db.DB.C("Prediction").Pipe(pipeline).All(&result)
// 	if err != nil {
// 		c.JSON(500, err)
// 		return
// 	}
// 	c.JSON(200, result)
// }
