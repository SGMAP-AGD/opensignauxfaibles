function finalize(k, v) {

    v = Object.keys((v.batch || {})).sort().reduce((m, batch) => {
        Object.keys(v.batch[batch]).forEach((type) => {
            m[type] = (m[type] || {})
            Object.assign(m[type], v.batch[batch][type])
        })
        return m
    }, { "siren": v.siren })

    liste_periodes = generatePeriodSerie(date_debut, date_fin)

    var value_array = liste_periodes.map(function (e) {
        return {
            "siren": v.siren,
            "periode": e,
            "arrete_bilan": new Date(0)
        }
    });

    Object.keys(v.bdf).forEach(hash => {
        value_array.forEach(periode => {
            if (periode.arrete_bilan.getTime() < v.bdf[hash].arrete_bilan.getTime()
                && periode.periode.getTime() > v.bdf[hash].arrete_bilan.getTime()) {
                periode.arrete_bilan = v.bdf[hash].arrete_bilan
                periode.secteur = v.bdf[hash].secteur
                periode.poids_frng = v.bdf[hash].poids_frng
                periode.taux_marge = v.bdf[hash].taux_marge
                periode.delai_fournisseur = v.bdf[hash].delai_fournisseur
                periode.dette_fiscale = v.bdf[hash].dette_fiscale
                periode.financier_court_terme = v.bdf[hash].financier_court_terme
                periode.frais_financier = v.bdf[hash].frais_financier
            }
        })
    })

    value_array.forEach((periode, index) => {
        if (periode.arrete_bilan.getTime() == 0) {
            delete value_array[index]
        }
    })

    return {"siren": k, "entreprise": value_array}
}