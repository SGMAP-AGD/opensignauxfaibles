package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Value Objet racine pour le stockage établissement.
// En insertion, ID est un ObjectId généré automatiquement, puis devient un siret après mapReduce
// Value est une map des lots d'intégration à ordonner par ordre alphanumérique
type Value struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
	// Value map[string]Etablissement `json:"value" bson:"value"`
	Value Etablissement `json:"value" bson:"value"`
}

// Etablissement objet établissement (/entreprise/)
type Etablissement struct {
	Siret             string             `json:"siret" bson:"siret"`
	AncienSiret       []string           `json:"ancien_siret" bson:"ancien_siret"`
	Compte            Compte             `json:"compte" bson:"compte"`
	Altares           map[string]Altares `json:"altares" bson:"altares"`
	ActivitePartielle ActivitePartielle  `json:"activite_partielle" bson:"activite_partielle"`
}

// Compte informations ursaff
type Compte struct {
	Effectif   map[string]Effectif   `json:"effectif" bson:"effectif"`
	Delais     map[string]Delais     `json:"delais" bson:"delais"`
	Debit      map[string]Debit      `json:"debit" bson:"debit"`
	CCSF       map[string]CCSF       `json:"ccsf" bson:"ccsf"`
	Cotisation map[string]Cotisation `json:"cotisation" bson:"cotisation"`
}

// Periode Période de temps avec un début et une fin
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}
