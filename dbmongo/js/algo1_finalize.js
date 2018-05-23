function finalize(k, v) {

    entreprise = Object.keys((v.batch||{})).sort().reduce((m,batch) => {
        Object.keys(v.batch[batch]).forEach((type) => {
            m[type] = (m[type] || {})
            Object.assign(m[type],v.batch[batch][type])
        })
        return m
    },{})

    etablissements = Object.keys((v.etablissement||{})).map(etablissement => {
        return Object.keys(v.etablissement[etablissement].batch).sort().reduce((m,batch) => {
            Object.keys(v.etablissement[etablissement].batch[batch]).forEach((type) => {
                m[type] = (m[type] || {})
                Object.assign(m[type],v.etablissement[etablissement].batch[batch][type])
            })
            return m
        }, {"entreprise": entreprise, "siret": etablissement})
    })

    return [].concat.apply([], etablissements.map(v => {

        v.apconso = (v.apconso||{})
        v.apdemande = (v.apdemande||{})
        v.effectif = (v.effectif || {})
        v.altares = (v.altares || {})
        v.cotisation = (v.cotisation || {})
        v.debit = (v.debit || {})

        // relier les débits
        var ecn = Object.keys(v.debit).reduce((m, h) => {
            var d = [h, v.debit[h]]
            var start = d[1].periode.start
            var end = d[1].periode.end
            var num_ecn = d[1].numero_ecart_negatif
            var compte = d[1].numero_compte
            var key = start + "-" + end + "-" + num_ecn + "-" + compte
            m[key] = (m[key] || []).concat([{
                "hash": d[0], 
                "numero_historique": d[1].numero_historique,
                "date_traitement": d[1].date_traitement
            }])
            return m 
        }, {})
        Object.keys(ecn).forEach(i => {

            ecn[i].sort(compareDebit)
            var l = ecn[i].length
            ecn[i].forEach((e,idx) => {

                if (idx <= l-2) {
                    v.debit[e.hash].debit_suivant = ecn[i][idx+1].hash;  
                }
            })
        })

        var date_debut = new Date("2014-01-01");
        var date_fin = new Date("2018-03-01");
        var date_fin_effectif = new Date("2018-01-01")

        var offset_effectif = (date_fin_effectif.getUTCFullYear()-date_fin.getUTCFullYear())*12 + date_fin_effectif.getUTCMonth()-date_fin.getUTCMonth()
        liste_periodes = generatePeriodSerie(date_debut, date_fin)

        //liste_periodes = [new Date('2015-01-01'),new Date('2016-01-01'),new Date('2018-02-01')]
        //liste_periodes = [new Date("2015-01-01"), new Date("2016-01-01"), new Date("2018-04-01")]

        var value_array = liste_periodes.map(function(e) {
            return {"siret": v.siret,
                    "periode": e, 
                    "lag_effectif_missing": true,
                    "apart_last12_months": false,
                    "apart_heures_consommees_array": [],
                    "effectif_history": {},
                    "outcome_0_12": "non-default",
                    "date_defaillance": null,
                    "cotisation_due_periode": {},
                    "debit_array": []
                }
        });


        var value = value_array.reduce(function(periode, val, index) {
            periode[val.periode.getTime()] = val     
            return periode
        }, {});

        map_effectif = Object.keys(v.effectif).reduce(function(map_effectif, hash) {
            var effectif = v.effectif[hash];
            if (effectif == null) {
                return map_effectif
            }
            var effectifTime = effectif.periode.getTime()
            map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
            return map_effectif
        }, {})
        

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



        // activite partielle    
        var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
            apart[v.apdemande[hash].id_demande] = {
                "demande": k,
                "consommation": []
            }
            return apart
        }, {})

        Object.keys(v.apconso).forEach(hash => {
            var valueap = v.apconso[k]
            if (valueap.id_conso.substring(0,10) in apart) {
                apart[valueap.id_conso.substring(0,10)].consommation.push(k)
            }
        })

        Object.keys(apart).forEach(k => {
            v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
            for (j in apart[k].consommation) {
                v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande;
            }
        })

        Object.keys(v.apconso).map(
            function(h) {
                var conso = v.apconso[h]
                if (conso.hash_demande && v.apdemande[conso.hash_demande].motif_recours_se != 3) {
                    
                    var currentTime = conso.periode.getTime()
                    var beforeTime = new Date(conso.periode.getTime()).setFullYear(conso.periode.getFullYear()+1)
                    var pastYearTimes = generatePeriodSerie(new Date(currentTime), new Date(beforeTime)).map(function(date) {return date.getTime()})
                    pastYearTimes.map(function(time) {
                        if (time in value) {
                            value[time].apart_last12_months = true;
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

        Object.keys(v.cotisation).map(function(h) {
            var cotisation = v.cotisation[h]
            var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
            periode_cotisation.map(function (date_cotisation) {
                value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
            })
        })

        var value_dette = {}

        Object.keys(v.debit).map(function(h) {
            var debit = v.debit[h]
            if (debit.part_ouvriere + debit.part_patronale > 0) {

                var debit_suivant = (v.debit[debit.debit_suivant] || debit)
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
            if (!(val.effectif)) {
                delete value[val.periode.getTime()]
                delete value_array[index]
            }
        })

        value_array.map(function(val,index) {

            val.lag_effectif_missing = (val.lag_effectif ? false : true)

            val.growthrate_effectif = (val.lag_effectif_missing ? 0 : val.effectif/val.lag_effectif)


            if (val.lag_effectif_missing) {
                val.cut_growthrate = "manquant"
            } else if (val.growthrate_effectif < 0.8) {
                val.cut_growthrate = "moins_de_20p"
            } else if (val.growthrate_effectif < 0.95) {
                val.cut_growthrate = "moins_20_a_5p" 
            } else if (val.growthrate_effectif < 1.05) {
                val.cut_growthrate = "stable"
            } else if (val.growthrate_effectif < 1.20) {
                val.cut_growthrate = "plus_5_a_20p"
            } else {
                val.cut_growthrate = "plus_20p"
            }

            if ((val.effectif || 0) <= 20) {
                val.cut_effectif = "10_20"
            } else if (val.effectif <= 50) {
                val.cut_effectif = "21_50"
            } else {
                val.cut_effectif = "Plus_de_50"
            }

            var e = Object.keys(val.effectif_history).reduce(function (m,h) {
                m.total += val.effectif_history[h];
                m.length += 1;
                return m
            },{"length": 0, "total": 0})

            val.effectif_average = e.total / e.length;

            val.apart_heures_consommees = val.apart_heures_consommees_array.reduce(function(m,h) {return m + h}, 0)
            val.apart_share_heuresconsommees = ((val.effectif_average||0) == 0 ? 0 : val.apart_heures_consommees / (val.effectif_average * 1607) * 100)

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
            
            val.indicatrice_dettecumulee_12m = (val.montant_part_ouvriere + val.montant_part_patronale) > 0
            
            val.ratio_dettecumulee_cotisation_12m = (val.mean_cotisation_due > 0 ? (val.montant_part_ouvriere + val.montant_part_patronale) / val.mean_cotisation_due : 0)
            val.log_ratio_dettecumulee_cotisation_12m = Math.log((val.ratio_dettecumulee_cotisation_12m + 1||1))
            val.apart_last12_months = (val.apart_last12_months?1:0)
            val.apart_consommee = (val.apart_heures_consommees>0?1:0)

            delete val.effectif_history
            delete val.cotisation_array
            delete val.debit_array
            delete val.montant_dette
            delete val.apart_heures_consommees_array
            delete val.cotisation_due_periode
            // delete val.date_defaillance
            delete val.montant_part_ouvriere
            delete val.montant_part_patronale
            delete val.ratio_dettecumulee_cotisation_12m
            delete val.mean_cotisation_due
            delete val.effectif_date
            delete val.effectif_average
            delete val.lag_effectif
        })

        return value_array
    }))
    
}


