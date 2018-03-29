function finalize(key, value) {

    
    v = value.value
    var periode = new Date('2017-05-01');

    var periode_before = new Date(periode.getTime())
    periode_before.setYear(periode_before.getYear()+1899);

    var periode_after = new Date(periode.getTime())
    periode_after.setYear(periode_after.getYear()+1901);

    var hash_effectif_all = Object.keys(v.compte.effectif).reduce(function(m, h) {
        var e = v.compte.effectif
        var compte = e[h].numero_compte
        if (e[h].periode.getTime() == periode.getTime()) {
                m.now[compte] = h}
        if (e[h].periode.getTime() == periode_before.getTime()) {
                m.before[compte] = h}
        if (periode_before <= e[h].periode && e[h].periode < periode) {
            m.past_year[e[h].periode] = (m.past_year[e[h].periode] || 0) + e[h].effectif
        }
        return m
    }, {"now": {}, "before": {}, "past_year": {}})

    var hash_effectif = hash_effectif_all.now
    var hash_effectif_before = hash_effectif_all.before

    var effectif = Object.keys(hash_effectif).reduce(function(m,h) {
        return m + v.compte.effectif[hash_effectif[h]].effectif
    }, 0)

    var lag_effectif = Object.keys(hash_effectif_before).reduce(function(m,h) {
        return m + v.compte.effectif[hash_effectif_before[h]].effectif
    }, 0)

    var lag_effectif_missing = (Object.keys(hash_effectif_before).length === 0)
    
    var growthrate_effectif = (lag_effectif_missing ? 0 : effectif/lag_effectif)

    if (lag_effectif_missing) {
        cut_growthrate = "manquant"
    } else if (growthrate_effectif < 0.8) {
        cut_growthrate = "moins de 20%"
    } else if (growthrate_effectif < 0.95) {
        cut_growthrate = "moins 20 à 5%" 
    } else if (growthrate_effectif < 1.05) {
        cut_growthrate = "stable"
    } else if (growthrate_effectif < 1.20) {
        cut_growthrate = "plus 5 à 20%"
    } else {
        cut_growthrate = "plus 20%"
    }

    if (effectif <= 20) {
        cut_effectif = "10-20"
    } else if (effectif <= 50) {
        cut_effectif = "21-50"
    } else {
        cut_effectif = "Plus de 50"
    }

    // Activité partielle
    var apart_last12_months = Object.keys(v.activite_partielle.demande).some(function(h) {
        if (v.activite_partielle.demande[h].periode.start >= periode_before && 
        v.activite_partielle.demande[h].periode.start < periode) {
            return true;
        }
    })

    var apart_heures_consommees = Object.keys(v.activite_partielle.consommation).reduce(function(m, h) {
        var c = v.activite_partielle.consommation[h]
        var d = (c.hash_demande ? v.activite_partielle.demande[c.hash_demande] : {"motif_recours_se": 1})
        return (periode_before < c.date && c.date <= periode && d.motif_recours_se != 3 ? m + c.heure_consomme : m)
    }, 0)

    var apart_consommee = (apart_heures_consommees > 0)

    var effectif_average = (Object.keys(hash_effectif_all.past_year).reduce(function(m,k) {
        return m + hash_effectif_all.past_year[k]
    }, 0) / Object.keys(hash_effectif_all.past_year).length || 0)

    var apart_share_heures_consommees = (effectif_average == 0 ? 0 : apart_heures_consommees / (effectif_average * 1607) * 100)

    // // defaillance - On prend la date de l'évènement le plus proche dans l'avenir par rapport à period
    var date_defaillance = Object.keys(v.altares).reduce(function(mem, hash) {
        return (isRJLJ(v.altares[hash].code_evenement) && 
                (mem || new Date()) > v.altares[hash].date_effet && 
                v.altares[hash].date_effet > periode ?
                v.altares[hash].date_effet : mem)
    }, null)

    var outcome_0_12 = (date_defaillance < periode_after && date_defaillance >= periode ? 'default' : 'non_default')

    // // Dette
    var cotisation_due = Object.keys(v.compte.cotisation).reduce(function(m,h) {
        var cotisation = v.compte.cotisation[h]

        if (cotisation.periode.end > periode_before && periode > cotisation.periode.start) {
            var start = new Date(Math.max(cotisation.periode.start, periode_before))

            var end = new Date(Math.min(cotisation.periode.end, periode))

            var duration_periode_cotisation = cotisation.periode.end.getMonth() - cotisation.periode.start.getMonth()
                + (12 * (cotisation.periode.end.getFullYear() - cotisation.periode.start.getFullYear()))

            var duration_periode = end.getMonth() - start.getMonth()
                + (12 * (end.getFullYear() - start.getFullYear()))
            p = {"start": start,
                 "end":  end}

            m.periode[cotisation.periode.start] = p
            m.du = m.du + cotisation.du * (duration_periode / duration_periode_cotisation)
            return m
        } else {
            return m
        }
    }, {"periode": {}, "du": 0})

    var periode_past = new Date(periode_before.getTime())
    var nb_month = 0
    while (periode_past < periode) {
        for (p in cotisation_due.periode) {
            if (cotisation_due.periode[p].start <= periode_past && periode_past < cotisation_due.periode[p].end) {
                nb_month += 1
                break
            }
        }
        periode_past.setMonth(periode_past.getMonth() + 1)
    }

    var mean_cotisation_due = cotisation_due.du / nb_month

    var log_cotisationdue_effectif = ( mean_cotisation_due * effectif == 0 ? 0 : Math.log(mean_cotisation_due / effectif))

    var dette = Object.keys(v.compte.debit).reduce(function(m,h) {
        var debit = v.compte.debit[h]
        var start = debit.periode.start.toISOString().substring(2,7).replace("-","")
        var end = debit.periode.end.toISOString().substring(2,7).replace("-","")
        var num_ecn = debit.numero_ecart_negatif
        var compte = debit.numero_compte
        var key = start + "-" + end + "-" + num_ecn + "-" + compte
        if ((!(key in m) || m[key].numero_historique < debit.numero_historique) && debit.date_traitement < periode) {
            m[key] = debit
        }
        return m
    }, {})
 
    var dette_cumulee = Object.keys(dette).reduce(function(m,k) {
        m.part_patronale += dette[k].part_patronale
        m.part_ouvriere += dette[k].part_ouvriere
        if (dette[k].periode.start >= periode_before) {
            m.part_patronale12m += dette[k].part_patronale
            m.part_ouvriere12m  += dette[k].part_ouvriere
        }
        return m
    }, {"part_patronale": 0,
        "part_ouvriere": 0,
        "part_patronale12m": 0,
        "part_ouvriere12m": 0}
    )

    var indicatrice_dettecumulee_12m = (dette.part_patronale12m + dette.part_ouvriere12m > 0)
    var ratio_dettecumulee_cotisation_12m = (dette_cumulee.part_ouvriere12m + dette_cumulee.part_patronale12m) / mean_cotisation_due
    var log_ratio_dettecumulee_cotisation_12m = (ratio_dettecumulee_cotisation_12m == 0 ? 0 : Math.log(ratio_dettecumulee_cotisation_12m))
    
    emit_value = {
        "siret": v.siret,
        "outcome_0_12": outcome_0_12,
        "cut_effectif": cut_effectif,
        "cut_growthrate": cut_growthrate,
        "lag_effectif_missing": lag_effectif_missing,
        "apart_last12months": apart_last12_months,
        "apart_consommee": apart_consommee,
        "apart_share_heuresconsommees": apart_share_heures_consommees,
        "log_cotisationdue_effectif": log_cotisationdue_effectif,
        "log_ratio_dettecumulee_cotisation_12m": log_ratio_dettecumulee_cotisation_12m,
        "indicatrice_dettecumulee_12m": indicatrice_dettecumulee_12m
    }

    return emit_value
    
}


