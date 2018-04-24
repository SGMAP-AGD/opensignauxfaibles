package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Value Objet racine pour le stockage établissement.
// En insertion, ID est un ObjectId généré automatiquement, puis devient un siret après mapReduce
// Value est une map des lots d'intégration à ordonner par ordre alphanumérique
type Value struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Etablissement `json:"value" bson:"value"`
}

// Etablissement objet établissement (/entreprise/)
type Etablissement struct {
	Siret       string           `json:"siret" bson:"siret"`
	Region      string           `json:"region" bson:"region"`
	Key         string           `json:"-" bson:"-"`
	AncienSiret []string         `json:"ancien_siret,omitempty" bson:"ancien_siret,omitempty"`
	Batch       map[string]Batch `json:"batch" bson:"batch"`
}

// Batch lot de data
type Batch struct {
	Effectif       map[string]Effectif       `json:"effectif,omitempty" bson:"effectif,omitempty"`
	Delais         map[string]Delais         `json:"delais,omitempty" bson:"delais,omitempty"`
	Debit          map[string]Debit          `json:"debit,omitempty" bson:"debit,omitempty"`
	CCSF           map[string]CCSF           `json:"ccsf,omitempty" bson:"ccsf,omitempty"`
	Cotisation     map[string]Cotisation     `json:"cotisation,omitempty" bson:"cotisation,omitempty"`
	Altares        map[string]Altares        `json:"altares,omitempty" bson:"altares,omitempty"`
	APDemande      map[string]APDemande      `json:"apdemande,omitempty" bson:"apdemande,omitempty"`
	APConsommation map[string]APConsommation `json:"apconso,omitempty" bson:"apconso,omitempty"`
}

// Periode Période de temps avec un début et une fin
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}
