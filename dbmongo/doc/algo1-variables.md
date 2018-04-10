# Algorithme 1, implémentation MongoDB

Traitement des variables

Cette documentattion a été établie à partir du commit `6534e6e`
Le code complet de la procédure est librement consultable ici: 
https://github.com/entrepreneur-interet-general/opensignauxfaibles/blob/6534e6e93d173af89a48b5c69581f6610d8eb6c2/dbmongo/finalize2.js

## Modèle de données
Les calculs sont effectués entreprise par entreprise. Voici un exemple abrégé des données telles qu'elles sont disposées dans le modèle:

```json
{"activite_partielle": {
    "consommation": {
        "1a2610c7be758ecc0014c42950a4b67c": {
            "effectif": 12,
            "hash_demande": "714b5be1bbc7cb5b338911c715df358b",
            "heure_consomme": 966,
            "id_conso": "02101280800",
            "montant": 7476.84,
            "periode": "2017-12-01T00:00:00Z"
        },
    [...]
    },
    "demande": {
        "38b4bf300d01aaf7a26284ec1e581823": {
            "avis_ce": 1,
            "date_statut": "2015-05-29T00:00:00Z",
            "effectif": 63,
            "effectif_autorise": 50,
            "effectif_consomme": 0,
            "effectif_entreprise": 170,
            "hash_consommation": [],
            "heure_consommee": 0,
            "hta": 40000,
            "id_demande": "0210128030",
            "montant_consommee": 0,
            "motif_recours_se": 1,
            "mta": 289200,
            "perimetre": 1,
            "periode": {
                "end": "2015-12-18T00:00:00Z",
                "start": "2015-06-19T00:00:00Z"
            },
            "prod_hta_effectif": 2000000,
            "recours_anterieur": 2,
            "tx_pc": 7.23,
            "tx_pc_etat_dares": 4.33,
            "tx_pc_unedic_dares": 2.9
        },
    [...]
    },
"altares": {},
"compte": {
    "ccsf": {},
    "cotisation": {
        "035d811d43a2d1efa59b4ac7c2c62ae4": {
            "du": 196186,
            "ecriture": "100",
            "encaisse": 196186,
            "numero_compte": "267000001620002812",
            "periode": {
                "end": "2015-12-01T00:00:00Z",
                "start": "2015-11-01T00:00:00Z"
            },
            "periode_debit": "115422000281100",
            "recouvrement": 0
            },
    [une clé pour chaque ligne]
    },
    "debit": {
        "0a3a8b6c5ae1009fc2e179a38a42ef1e": {
            "code_motif_ecart_negatif": "19",
            "code_operation_ecart_negatif": "11",
            "code_procedure_collective": "0",
            "date_traitement": "2012-02-08T00:00:00Z",
            "debit_suivant": "83ee6b542a88fb3f285e3873269c04be",
            "etat_compte": 1,
            "numero_compte": "267000001620002812",
            "numero_ecart_negatif": "101",
            "numero_historique": 3,
            "part_ouvriere": 0,
            "part_patronale": 654,
            "periode": {
                "end": "2012-01-01T00:00:00Z",
                "start": "2011-01-01T00:00:00Z"
            }
        },
    [idem]
    },
    "delais": {},
    "effectif": {
        "001f36e097c86a7b0965d3ac3ebf1dd5": {
            "effectif": 122,
            "numero_compte": "267000001620002812",
            "periode": "2013-05-01T00:00:00Z"
            },
    [...]
    },
},
"index": {
    "algo1": true
    },
"siret": "01545141200025"
}
```

Cette structure de données est identifiée par la variable `v` dans le code.

Chaque ligne d'information est contenue dans une clé hashée qui permet l'identification et la comparaison rapide dans la gestion des lots d'intégration mois après mois. Cette clé est représentée par la variable `h`.

Le traitement construit et retourne la variable `array_value` dont les éléments sont également accessibles par l'objet `value` sous forme d'un dictionnaire des périodes temporelles.

Lorsque des traitements itèrent sur les périodes d'`array_value`, la variable utilisée pour accéder à l'élément est `val`.

## Cotisations (`mean_cotisation_due`)

La traitement constitue un calendrier des cotisations en remettant sur une base mensuelle toutes les sommes de la base.

```js
var value_cotisation = {}

Object.keys(v.compte.cotisation).map(function(h) {
    var cotisation = v.compte.cotisation[h]
    var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.map(function (date_cotisation) {
        value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
})
```

Ensuite la procédure rattache les périodes de ce calendrier aux périodes de calcul du résultat (12 période prédécentes rattachées à la période en cours).

```js
// Itération sur toutes les clés de l'objet value
Object.keys(value).map(function(time) {

    // période actuelle
    var currentTime = value[time].periode.getTime()

    // première période passée
    var beforeTime = new Date(value[time].periode.getTime()).setFullYear(value[time].periode.getFullYear()-1)

    // génération de la série des 12 périodes précédentes
    var pastYearTimes = generatePeriodSerie(new Date(beforeTime), new Date(currentTime)).map(function(date) {return date.getTime()})

    // tableau d'accès aux 12 dernières périodes de cotisations
    value[time].cotisation_array = pastYearTimes.map(function(t) {
        return value_cotisation[t]
    })
})
```

Voici un exemple des données stockées dans cette variable temporaire pour une période donnée. Nous avons ici un aperçu des 2 premières périodes composées de 3 sources différentes (3 périodes ou comptes se chevauchant).

```js
{"cotisation_array": [
          [
            1271.5833333333333,
            86795,
            2.8333333333333335
          ],
          [
            120225,
            1271.5833333333333,
            2.8333333333333335
          ],
          [
            118042,
            1271.5833333333333,
            2.8333333333333335
          ],
    [...]
}
```

Pour finir, pour chaque période du calcul, on calcul la moyenne mensuelle de la cotisation des 12 derniers mois accessibles dans `val.cotisation_array`.

```js
value_array.map(function(val,index) {
    c = val.cotisation_array.reduce(function(m,cot) {
        m.nb_month += 1
        m.sum += (cot||[]).reduce((a,b) => a+b, 0)
        return m
    }, {"sum": 0, "nb_month": 0})
    val.mean_cotisation_due = (c.nb_month > 0 ? c.sum / c.nb_month : 0)
}
```


## Débit (`montant_part_patronale` et `montant_part_ouvriere`)

Lors du calcul des variables, la procédure constitue un calendrier mensuel et rattache les débits de l'entreprise sur les périodes concernées par les intervales de dates de traitement d'écart négatif (limité à 12 mois pour correspondre aux variables calculées dans le modèle R).

```js
var value_dette = {}

Object.keys(v.compte.debit).map(function(h) {
    var debit = v.compte.debit[h]
    if (debit.part_ouvriere + debit.part_patronale > 0) {
        var debit_suivant = (v.compte.debit[debit.debit_suivant] || debit)

        // Calcul de la date limite de prise en compte du débit (+12mois)
        date_limite = new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))

        // Date arrondie au mois supérieur 
        date_traitement_debut = new Date(
            Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + (debit.date_traitement.getUTCDate() > 1 ? 1:0))
        )

        // Date arrondie au mois supérieur
        date_traitement_fin = new Date(
            Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + (debit_suivant.date_traitement.getUTCDate() > 1 ? 1:0))
        )

        // On comptabilise un débit au maximum 12 mois après sa survenue
        periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)

        periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)

        // Pour toutes les périodes on écrit le débit dans le résultat
        generatePeriodSerie(periode_debut, periode_fin).map(function(date) {
            time = date.getTime()
            value_dette[time] = (value_dette[time] || []).concat([{"periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale}])
        })
    }
})
```

Voici un exemple de la valeur de `value_dette`. On retrouve dans chaque `time` tous les débits qui ont encore une valeur à régler à `Date(time)`.

```
 "value_dette": {
          "1443657600000": [ 
            {
              "part_ouvriere": 1,
              "part_patronale": 176,
              "periode": "2015-01-01T00:00:00Z"
            },
            {
              "part_ouvriere": 0,
              "part_patronale": 185,
              "periode": "2015-02-01T00:00:00Z"
            },
            {
              "part_ouvriere": 1,
              "part_patronale": 195,
              "periode": "2015-03-01T00:00:00Z"
            },
            [...]
          ]
    [...]
 }
```


Une fois ce calendrier constitué, les débits sont rattachés aux périodes de `value_array` en passant par les clés de l'object `value` qui jouent ici le rôle de pointeur.

```js
Object.keys(value).map(function(time) {
    var currentTime = value[time].periode.getTime()
    var beforeTime = new Date(value[time].periode.getTime()).setFullYear(value[time].periode.getFullYear()-1)
    var pastYearTimes = generatePeriodSerie(new Date(beforeTime), new Date(currentTime)).map(function(date) {return date.getTime()})

    //
    value[time].cotisation_array = pastYearTimes.map(function(t) {
        return value_cotisation[t]
    })
    if (time in value_dette) {
        value[time].debit_array = value_dette[time]
    }
})
```

Pour finir, le total est calculé pour chaque période:

```js
val.montant_dette = val.debit_array.reduce(function(m,dette) {
    m.part_ouvriere += dette.part_ouvriere
    m.part_patronale += dette.part_patronale
    return m
}, {"part_ouvriere": 0, "part_patronale": 0})
```
## Effectifs

Pour prendre en compte le décalage des effectifs, nous calculons un nombre de mois de décalage que l'on applique à toutes les affectations d'effectifs.

```js
var offset_effectif = (date_fin_effectif.getUTCFullYear()-date_fin.getUTCFullYear())*12 + date_fin_effectif.getUTCMonth()-date_fin.getUTCMonth()
```

On crée un dictionnaire des périodes d'effectifs accessible avec le timestamp d'une période.
```js
map_effectif = Object.keys(v.compte.effectif).reduce(function(map_effectif, hash) {
    var effectif = v.compte.effectif[hash];
    var effectifTime = effectif.periode.getTime()
    map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
    return map_effectif
}, {})
```

On enregistre dans `value_array` les données d'effectif sous 3 formes:

- l'effectif en cours (+ décalage = ``effectifDate`)
- l'effectif 12 mois avant (+ décalage = `historyDate`)
- les effectifs de 12 mois à maintenant (`historyPeriods`)

```js
    // inscription des effectifs dans les périodes
    value_array.map(function(val) {
        var currentTime = val.periode.getTime()
        var effectifDate = DateAddMonth(val.periode, offset_effectif)
        var historyDate = DateAddMonth(val.periode, offset_effectif - 12)
        var historyPeriods = generatePeriodSerie(historyDate, effectifDate)
        val.effectif_date = effectifDate
        val.effectif = map_effectif[effectifDate.getTime()]

        val.lag_effectif = map_effectif[historyDate.getTime()]
        historyPeriods.map(function(p) {
            val.effectif_history[p.getTime()] = map_effectif[p.getTime()]
        })
    })
```

## Variables produites

### `lag_effectif_missing`
```js
val.lag_effectif_missing = (val.lag_effectif ? false : true)
```

### `cut_effectif`
```js
if ((val.effectif || 0) <= 20) {
    val.cut_effectif = "10-20"
} else if (val.effectif <= 50) {
    val.cut_effectif = "21-50"
} else {
    val.cut_effectif = "Plus de 50"
}
```

### `cut_growthrate`
```js
if (val.lag_effectif_missing) {
    val.cut_growthrate = "manquant"
} else if (val.growthrate_effectif < 0.8) {
    val.cut_growthrate = "moins de 20%"
} else if (val.growthrate_effectif < 0.95) {
    val.cut_growthrate = "moins 20 à 5%" 
} else if (val.growthrate_effectif < 1.05) {
    val.cut_growthrate = "stable"
} else if (val.growthrate_effectif < 1.20) {
    val.cut_growthrate = "plus 5 à 20%"
} else {
    val.cut_growthrate = "plus 20%"
}
```

### ```log_cotisationdue_effectif```
```js
val.log_cotisationdue_effectif = (val.mean_cotisation_due * val.effectif == 0 ? 0 : Math.log(1+val.mean_cotisation_due/val.effectif))
```

### `log_ratio_dettecumulee_cotisation_12m`
```js
val.ratio_dettecumulee_cotisation_12m = (val.mean_cotisation_due > 0 ? (val.montant_part_ouvriere + val.montant_part_patronale) / val.mean_cotisation_due : 0)
    val.log_ratio_dettecumulee_cotisation_12m = Math.log((val.ratio_dettecumulee_cotisation_12m + 1||1))
```
### `indicatrice_dettecumulee_12m`
```js
val.indicatrice_dettecumulee_12m = (val.montant_part_ouvriere + val.montant_part_patronale) > 0
```

## Exemple de résultat
Le calcul est effectué pour les périodes `2015-01-01`, `2016-01-01` et `2018-02-01`

```js
[
  {
    "id": "01545141200025",
    "value": [
      {
        "apart_consommee": 0,
        "apart_heures_consommees": 0,
        "apart_last12_months": 0,
        "apart_share_heuresconsommees": 0,
        "cut_effectif": "Plus de 50",
        "cut_growthrate": "moins 20 à 5%",
        "effectif": 93,
        "growthrate_effectif": 0.808695652173913,
        "indicatrice_dettecumulee_12m": false,
        "lag_effectif_missing": false,
        "log_cotisationdue_effectif": 7.232944179437735,
        "log_ratio_dettecumulee_cotisation_12m": 0,
        "outcome_0_12": "non-default",
        "periode": "2015-01-01T00:00:00Z",
        "siret": "01545141200025"
      },
      {
        "apart_consommee": 0,
        "apart_heures_consommees": 0,
        "apart_last12_months": 0,
        "apart_share_heuresconsommees": 0,
        "cut_effectif": "Plus de 50",
        "cut_growthrate": "moins 20 à 5%",
        "effectif": 84,
        "growthrate_effectif": 0.9032258064516129,
        "indicatrice_dettecumulee_12m": false,
        "lag_effectif_missing": false,
        "log_cotisationdue_effectif": 7.2268473344668465,
        "log_ratio_dettecumulee_cotisation_12m": 0,
        "outcome_0_12": "non-default",
        "periode": "2016-01-01T00:00:00Z",
        "siret": "01545141200025"
      },
      {
        "apart_consommee": 1,
        "apart_heures_consommees": 9632,
        "apart_last12_months": 1,
        "apart_share_heuresconsommees": 9.812459303643422,
        "cut_effectif": "Plus de 50",
        "cut_growthrate": "stable",
        "effectif": 60,
        "growthrate_effectif": 0.9836065573770492,
        "indicatrice_dettecumulee_12m": false,
        "lag_effectif_missing": false,
        "log_cotisationdue_effectif": 7.071207800325173,
        "log_ratio_dettecumulee_cotisation_12m": 0,
        "outcome_0_12": "non-default",
        "periode": "2018-02-01T00:00:00Z",
        "siret": "01545141200025"
      }
    ]
  }
]
```