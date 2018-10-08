package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func etablissementBrowse(c *gin.Context) {

	batchKey := c.Params.ByName("batchKey")
	batch := getBatch(batchKey)

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	siret := c.Params.ByName("siret")
	if siret == "" {
		c.JSON(500, "Aucun siret dans la requête")
		return
	}
	query := bson.M{"value.siret": siret}

	// Ressources JS
	MREtablissement := MapReduceJS{}
	errEt := MREtablissement.load("browse", "etablissement")
	if errEt != nil {
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	naf, errNAF := loadNAF()
	if errNAF != nil {
		c.JSON(500, "Problème d'accès aux fichiers naf")
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      nil,
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

	var result interface{}

	_, err := db.DB.C("Etablissement").Find(query).MapReduce(job, &result)

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, result)
	}

}
func predictionBrowse(c *gin.Context) {
	var result []interface{}
	db.DB.C("Prediction").Find(nil).All(&result)
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
