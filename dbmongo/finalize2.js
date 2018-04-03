function finalize(k, v) {

    var date_debut = new Date("2014-01-01");
    var date_fin = new Date("2018-04-01");

    var value_array = generatePeriodSerie(date_debut, date_fin).map(function(e) {
        return {"periode": e, 
                "lag_effectif_missing": true,
                "apart_last12months": false,
                "apart_heures_consommees_array": [],
                "past_year_effectif": {},
                "outcome_0_12": "non_default",
                "date_defaillance": null,
                "cotisation_due_periode": {},
                "debit_array": []
            }
    });

    var value = value_array.reduce(function(periode, val, index) {
        periode[val.periode.getTime()] = val     
        return periode
    }, {});



    Object.keys(v.compte.effectif).map(function(hash) {
        var e = v.compte.effectif[hash];
        var currentTime = e.periode.getTime()
        var beforeTime = (new Date(currentTime)).setFullYear(e.periode.getFullYear() + 1)
        var pastYearTimes = generatePeriodSerie(new Date(currentTime), new Date(beforeTime)).map(function(date) {return date.getTime()})
     
        pastYearTimes.map(function (time) {
            if (time in value) {
                value[time].past_year_effectif[currentTime] = (value[time].past_year_effectif[currentTime] || 0) + e.effectif 
            }
        })
        
        if (currentTime in value) {
            val = value[currentTime]
            val.effectif = (val.effectif || 0) + e.effectif;
        }

        if (beforeTime in value) {
            val = value[beforeTime]
            val.lag_effectif = (val.lag_effectif || 0) + e.effectif;
            val.lag_effectif_missing = false;
        }
    })

    value_array.map(function(val,index) {
        if (!(val.effectif)) {
            delete value[val.periode.getTime()]
            delete value_array[index]
        }
    })

    // activite partielle

    Object.keys(v.activite_partielle.consommation).map(
        function(h) {
            var conso = v.activite_partielle.consommation[h]
            if (conso.hash_demande && v.activite_partielle.demande[conso.hash_demande].motif_recours_se != 3) {
                
                var currentTime = conso.periode.getTime()
                var beforeTime = new Date(conso.periode.getTime()).setFullYear(conso.periode.getFullYear()+1)
                var pastYearTimes = generatePeriodSerie(new Date(currentTime), new Date(beforeTime)).map(function(date) {return date.getTime()})
                pastYearTimes.map(function(time) {
                    if (time in value) {
                        value[time].apart_last12months = true;
                        value[time].apart_heures_consommees_array.push(conso.heure_consomme);
                    }
                })     
            }
        })
    

    // defaillance - On prend la date de l'évènement le plus proche dans l'avenir par rapport à period
    Object.keys(v.altares).map(
        function(hash) {
            var altares = v.altares[hash]
            var periode_effet = new Date(Date.UTC(altares.date_effet.getUTCFullYear(), altares.date_effet.getUTCMonth(), 1, 0, 0, 0, 0))
            var periode_outcome = new Date(Date.UTC(altares.date_effet.getUTCFullYear() - 1, altares.date_effet.getUTCMonth(), 1, 0, 0, 0, 0))
            var pastYearTimes = generatePeriodSerie(periode_outcome, periode_effet).map(function(date) {return date.getTime()})
            pastYearTimes.map(
                function(time) {
                    if (time in value) {
                        value[time].date_defaillance = altares.date_effet
                        value[time].outcome_0_12 = "default";
                    }
                }
            )
        }
    )

    // Cotisation
    var value_cotisation = {}

    Object.keys(v.compte.cotisation).map(function(h) {
        var cotisation = v.compte.cotisation[h]
        var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
        periode_cotisation.map(function (date_cotisation) {
            value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
        })
    })




    var value_dette = {}

    Object.keys(v.compte.debit).map(function(h) {
        var debit = v.compte.debit[h]
        if (debit.part_ouvriere + debit.part_patronale > 0) {

            var debit_suivant = (v.compte.debit[debit.debit_suivant] || debit)
            date_limite = new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
            date_traitement_debut = new Date(
                Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + (debit.date_traitement.getUTCDate() > 1 ? 1:0))
            )

            date_traitement_fin = new Date(
                Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + (debit_suivant.date_traitement.getUTCDate() > 1 ? 1:0))
            )

            periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
            periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)

            generatePeriodSerie(periode_debut, periode_fin).map(function(date) {
                time = date.getTime()
                value_dette[time] = (value_dette[time] || []).concat([{"periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale}])
            })
        }
    })
 
    Object.keys(value).map(function(time) {
        var currentTime = value[time].periode.getTime()
        var beforeTime = new Date(value[time].periode.getTime()).setFullYear(value[time].periode.getFullYear()-1)
        var pastYearTimes = generatePeriodSerie(new Date(beforeTime), new Date(currentTime)).map(function(date) {return date.getTime()})
        value[time].cotisation_array = pastYearTimes.map(function(t) {
            return value_cotisation[t]
        })
        if (time in value_dette) {
            value[time].debit_array = value_dette[time]
        }
    })



    value_array.map(function(val,index) {
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

        if ((val.effectif || 0) <= 20) {
            val.cut_effectif = "10-20"
        } else if (val.effectif <= 50) {
            val.cut_effectif = "21-50"
        } else {
            val.cut_effectif = "Plus de 50"
        }

        val.growthrate_effectif = (val.lag_effectif_missing ? 0 : val.effectif/val.lag_effectif)
        var e = Object.keys(val.past_year_effectif).reduce(function (m,h) {
            m.total += val.past_year_effectif[h];
            m.length += 1;
            return m
        },{"length": 0, "total": 0})

        val.effectif_average = e.total / e.length;

        val.apart_heures_consommees = val.apart_heures_consommees_array.reduce(function(m,h) {return m + h}, 0)
        val.apart_share_heures_consommees = (val.effectif_average == 0 ? 0 : val.apart_heures_consommees / (val.effectif_average * 1607) * 100)

        c = val.cotisation_array.reduce(function(m,cot) {
            m.nb_month += 1
            m.sum += (cot||[]).reduce((a,b) => a+b, 0)
            return m
        }, {"sum": 0, "nb_month": 0})
        val.mean_cotisation_due = (c.nb_month > 0 ? c.sum / c.nb_month : 0)

        val.log_cotisationdue_effectif = (val.mean_cotisation_due * val.effectif == 0 ? 0 : Math.log(1+val.mean_cotisation_due/val.effectif))

        val.montant_dette = val.debit_array.reduce(function(m,dette) {
            m.part_ouvriere += dette.part_ouvriere
            m.part_patronale += dette.part_patronale
            return m
        }, {"part_ouvriere": 0, "part_patronale": 0})

        val.montant_part_ouvriere = val.montant_dette.part_ouvriere
        val.montant_part_patronale = val.montant_dette.part_patronale

        indicatrice_dettecumulee = (val.montant_part_ouvriere + val.montant_part_patronale) > 0

        val.ratio_dettecumulee_cotisation = (val.mean_cotisation_due > 0 ? (val.montant_part_ouvriere + val.montant_part_patronale) / val.mean_cotisation_due : 0)
        val.log_ratio_dettecumulee_cotisation = Math.log((val.ratio_dettecumulee_cotisation + 1||1))
        delete val.past_year_effectif
        delete val.cotisation_array
        delete val.debit_array
        delete val.montant_dette
        delete val.apart_heures_consommees_array

        if (!(val.effectif)) {delete value_array[index]}
    })
    
    return value_array
    
}


