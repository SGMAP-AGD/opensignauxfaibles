function finalize(k, v) {
    var entreprise = (v.entreprise||[]).reduce((m, periode) => {
        if (periode) {
            m[periode.periode.getTime()] = periode
        }
        return m
    }, {})
    empty = {
        "arrete_bilan": null,
        "secteur": null,
        "poids_frng": null,
        "taux_marge": null,
        "delai_fournisseur": null,
        "dette_fiscale": null,
        "financier_court_terme": null,
        "frais_financier": null
    }
    return Object.keys(v).reduce((m, k) => {
        if (k != 'entreprise' && k != "siren") {

            v[k].forEach(periode => {
                if (periode) {
                    Object.assign(periode, (entreprise[(periode || { "periode": new Date(0) }).periode.getTime()] || empty))
                    m.push(periode)
                }
            })

        }
        return m

    }, [])
}