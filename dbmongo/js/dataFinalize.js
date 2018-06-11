function finalize(k, v) {

    var result = {}
    var offset_effectif = (date_fin_effectif.getUTCFullYear() - date_fin.getUTCFullYear()) * 12 + date_fin_effectif.getUTCMonth() - date_fin.getUTCMonth()
    var liste_periodes = serie_periode

    v.apconso = (v.apconso || {})
    v.apdemande = (v.apdemande || {})
    v.effectif = (v.effectif || {})
    v.altares = (v.altares || {})
    v.cotisation = (v.cotisation || {})
    v.debit = (v.debit || {})
    var siret = v.siret
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
        ecn[i].forEach((e, idx) => {
            if (idx <= l - 2) {
                v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
            }
        })
    })

    var value_array = liste_periodes.map(function (e) {
        return {
            "siret": v.siret,
            "periode": e,
            "effectif_history": {},
            "cotisation_due_periode": {},
            "debit_array": []
        }
    });

    var value = value_array.reduce(function (periode, val, index) {
        periode[val.periode.getTime()] = val
        return periode
    }, {});

    map_effectif = Object.keys(v.effectif).reduce(function (map_effectif, hash) {
        var effectif = v.effectif[hash];
        if (effectif == null) {
            return map_effectif
        }
        var effectifTime = effectif.periode.toISOString()
        map_effectif[effectif.periode.toISOString()] = (map_effectif[effectif.periode.toISOString()] || 0) + effectif.effectif
        return map_effectif
    }, {})

    // inscription des effectifs dans les périodes
    value_array.map(function (val) {
        var currentTime = val.periode.getTime()
        var effectifDate = DateAddMonth(val.periode, offset_effectif)
        var historyDate = DateAddMonth(val.periode, offset_effectif - 12)
        var historyPeriods = generatePeriodSerie(historyDate, effectifDate)
        val.effectif_date = effectifDate
        val.effectif = map_effectif[effectifDate.getTime()]
        val.lag_effectif = map_effectif[historyDate.getTime()]

    })

    result.map_effectif = map_effectif
    // Cotisation
    var value_cotisation = {}

    Object.keys(v.cotisation).map(function (h) {
        var cotisation = v.cotisation[h]
        var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
        periode_cotisation.map(function (date_cotisation) {
            value_cotisation[date_cotisation.toISOString()] =
                (value_cotisation[date_cotisation.toISOString()] || []).concat(
                    {
                        "montant_cotisation": cotisation.du / periode_cotisation.length,
                        "ecriture": cotisation.ecriture,
                        "numero_compte": cotisation.numero_compte,
                        "periode_debit": cotisation.periode_debit,
                        "cotisation_totale": cotisation.du,
                        "duree_periode": periode_cotisation.length,
                        "periode_debut_cotisation": cotisation.periode.start,
                        "periode_fin_cotisation": cotisation.periode.end,
                        "periode_algorithme": date_cotisation.toISOString()
                    }
                )
        })
    })

    result.value_cotisation = value_cotisation

    var value_dette = {}

    Object.keys(v.debit).forEach(function (h) {
        var debit = v.debit[h]
        if (debit.part_ouvriere + debit.part_patronale > 0) {

            var debit_suivant = (v.debit[debit.debit_suivant] || debit)
            date_limite = new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
            date_traitement_debut = new Date(
                Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + (debit.date_traitement.getUTCDate() > 1 ? 1 : 0))
            )

            date_traitement_fin = new Date(
                Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + (debit_suivant.date_traitement.getUTCDate() > 1 ? 1 : 0))
            )

            periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
            periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)

            generatePeriodSerie(periode_debut, periode_fin).map(function (date) {
                time = date.getTime()
                value_dette[date.toISOString()] = (value_dette[date.toISOString()] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale }])
            })
        }
    })

    result.value_dette = value_dette
    result.siret = siret
    return result
}