package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Siret Système Informatique pour le Répertoire des Entreprises sur le Territoire
type Siret string

// Siret instancie une variable de type Siret
func (s *Siret) Siret() Siret {
	return "test"
}

// Value mongodb mapReduce object like
type Value struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
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

// ActivitePartielle Informations d'activité partielle
type ActivitePartielle struct {
	Demande      map[string]APDemande      `json:"demande" bson:"demande"`
	Consommation map[string]APConsommation `json:"consommation" bson:"consommation"`
}

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `json:"id_demande" bson:"id_demande"`
	EffectifEntreprise int       `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           int       `json:"effectif" bson:"effectif"`
	DateStatut         time.Time `json:"date_statut" bson:"date_statut"`
	TxPC               float64   `json:"tx_pc" bson:"tx_pc"`
	TxPCUnedicDares    float64   `json:"tx_pc_unedic_dares" bson:"tx_pc_unedic_dares"`
	TxPCEtatDares      float64   `json:"tx_pc_etat_dares" bson:"tx_pc_etat_dares"`
	Periode            Periode   `json:"periode" bson:"periode"`
	HTA                float64   `json:"hta" bson:"hta"`
	MTA                float64   `json:"mta" bson:"mta"`
	EffectifAutorise   int       `json:"effectif_autorise" bson:"effectif_autorise"`
	ProdHTAEffectif    float64   `json:"prod_hta_effectif" bson:"prod_hta_effectif"`
	MotifRecoursSE     int       `json:"motif_recours_se" bson:"motif_recours_se"`
	Perimetre          int       `json:"perimetre" bson:"perimetre"`
	RecoursAnterieur   int       `json:"recours_anterieur" bson:"recours_anterieur"`
	AvisCE             int       `json:"avis_ce" bson:"avis_ce"`
	HeureConsommee     float64   `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    float64   `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   int       `json:"effectif_consomme" bson:"effectif_consomme"`
}

// APConsommation Consommation d'activité partielle
type APConsommation struct {
	ID             string    `json:"id_conso" bson:"id_conso"`
	HeureConsommee float64   `json:"heure_consomme" bson:"heure_consomme"`
	Montant        float64   `json:"montant" bson:"montant"`
	Effectif       int       `json:"effectif" bson:"effectif"`
	Date           time.Time `json:"date" bson:"date"`
}

// Compte informations ursaff
type Compte struct {
	Effectif   map[string]Effectif   `json:"effectif" bson:"effectif"`
	Delais     map[string]Delais     `json:"delais" bson:"delais"`
	Debit      map[string]Debit      `json:"debit" bson:"debit"`
	CCSF       map[string]CCSF       `json:"ccsf" bson:"ccsf"`
	Cotisation map[string]Cotisation `json:"cotisation" bson:"cotisation"`
}

// Effectif Urssaf
type Effectif struct {
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
}

// CCSF information urssaf ccsf
type CCSF struct {
	DateTraitement time.Time `json:"date_traitement" bson:"date_traitement"`
	Stade          string    `json:"stade" bson:"stade"`
	Action         string    `json:"action" json:"action"`
}

// Delais tuple fichier ursaff
type Delais struct {
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheanche     time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	NumeroStructure   string    `json:"numero_structure" bson:"numero_structure"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

// Debit Débit – fichier Urssaf
type Debit struct {
	NumeroCompte                 string    `json:"numero_compte" bson:"numero_compte"`
	NumeroEcartNegatif           string    `json:"numero_ecart_negatif" bson:"numero_ecart_negatif"`
	DateTraitement               time.Time `json:"date_traitement" bson:"date_traitement"`
	PartOuvriere                 float64   `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale                float64   `json:"part_patronale" bson:"part_patronale"`
	NumeroHistoriqueEcartNegatif int       `json:"numero_historique" bson:"numero_historique"`
	EtatCompte                   int       `json:"etat_compte" bson:"etat_compte"`
	CodeProcedureCollective      string    `json:"code_procedure_collective" bson:"code_procedure_collective"`
	Periode                      Periode   `json:"periode" bson:"periode"`
	CodeOperationEcartNegatif    string    `json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
	CodeMotifEcartNegatif        string    `json:"code_motif_ecart_negatif" bson:"code_motif_ecart_negatif"`
}

// Cotisation Cotisation – fichier Urssaf
type Cotisation struct {
	NumeroCompte string  `json:"numero_compte" bson:"numero_compte"`
	PeriodeDebit string  `json:"periode_debit" bson:"periode_debit"`
	Periode      Periode `json:"period" bson:"periode"`
	Recouvrement float64 `json:"recouvrement" bson:"recouvrement"`
	Encaisse     float64 `json:"encaisse" bson:"encaisse"`
	Du           float64 `json:"du" bson:"du"`
	Ecriture     string  `json:"ecriture" bson:"ecriture"`
}

// Altares Extrait du récapitulatif altarès
type Altares struct {
	DateEffet     time.Time `json:"date_effet" bson:"date_effet"`
	DateParution  time.Time `json:"date_parution" bson:"date_parution"`
	CodeJournal   string    `json:"code_journal" bson:"code_journal"`
	CodeEvenement string    `json:"code_evenement" bson:"code_evenement"`
}

// Periode Période de temps avec un début et une fin
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}
