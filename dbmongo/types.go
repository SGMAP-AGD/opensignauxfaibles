package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Value Objet racine pour le stockage établissement.
// En insertion, ID est un ObjectId généré automatiquement, puis devient un siret après mapReduce
// Value est une map des lots d'intégration à ordonner par ordre alphanumérique
// type Value struct {
// 	ID    bson.ObjectId `json:"id" bson:"_id"`
// 	Value Etablissement `json:"value" bson:"value"`
// }

// ValueEntreprise permet de stocker une entreprise dans un objet Bson
type ValueEntreprise struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Entreprise    `json:"value" bson:"value"`
}

// ValueEtablissement structure pour un établissement
type ValueEtablissement struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Etablissement `json:"value" bson:"value"`
}

// Etablissement objet établissement (/entreprise/)
type Etablissement struct {
	Siret       string           `json:"siret" bson:"siret"`
	Region      string           `json:"region,omitempty" bson:"region,omitempty"`
	Key         string           `json:"-" bson:"-"`
	AncienSiret []string         `json:"ancien_siret,omitempty" bson:"ancien_siret,omitempty"`
	Batch       map[string]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// Entreprise object Entreprise
type Entreprise struct {
	Siren string           `json:"siren" bson:"siren"`
	Key   string           `json:"-" bson:"-"`
	Batch map[string]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// Batch lot de data
type Batch struct {
	Compact    map[string]bool       `json:"compact,omitempty" bson:"compact,omitempty"`
	Effectif   map[string]Effectif   `json:"effectif,omitempty" bson:"effectif,omitempty"`
	Delai      map[string]Delai      `json:"delai,omitempty" bson:"delai,omitempty"`
	Debit      map[string]Debit      `json:"debit,omitempty" bson:"debit,omitempty"`
	CCSF       map[string]CCSF       `json:"ccsf,omitempty" bson:"ccsf,omitempty"`
	Cotisation map[string]Cotisation `json:"cotisation,omitempty" bson:"cotisation,omitempty"`
	Altares    map[string]Altares    `json:"altares,omitempty" bson:"altares,omitempty"`
	APDemande  map[string]APDemande  `json:"apdemande,omitempty" bson:"apdemande,omitempty"`
	APConso    map[string]APConso    `json:"apconso,omitempty" bson:"apconso,omitempty"`
	Sirene     map[string]Sirene     `json:"sirene,omitempty" bson:"sirene,omitempty"`
	BDF        map[string]BDF        `json:"bdf,omitempty" bson:"bdf,omitempty"`
	Prediction map[string]Prediction `json:"prediction,omitempty" bson:"prediction,omitempty"`
	DPAE       map[string]DPAE       `json:"dpae,omitempty" bson:"dpae,omitempty"`
}

// Periode Période de temps avec un début et une fin
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}
