function map() { 
    v = Object.keys((this.value.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
        Object.keys(this.value.batch[batch])
            .filter(type => ['bdf', 'diane']
            .includes(type))
            .forEach((type) => {
                m[type] = (m[type] || {})
                var  array_delete = (this.value.batch[batch].compact.delete[type]||[])
                if (array_delete != {}) {array_delete.forEach(hash => {
                    delete m[type][hash]
                })
                }
                Object.assign(m[type], this.value.batch[batch][type])
        })
        return m
    }, { "siren": this.value.siren })

    v.bdf = Object.keys(v.bdf || {}).reduce((accu, key) => {
        var bdf = {
            annee: v.bdf[key].annee,
            arrete_bilan: v.bdf[key].arrete_bilan,
            poids_frng: v.bdf[key].poids_frng,
            taux_marge: v.bdf[key].taux_marge,
            delai_fournisseur: v.bdf[key].delai_fournisseur,
            dette_fiscale: v.bdf[key].dette_fiscale,
            financier_court_terme: v.bdf[key].financier_court_terme,
            frais_financier: v.bdf[key].frais_financier
        }
        accu = accu.concat(bdf)
        return accu
    }, []).sort((a1, a2) => a1.annee < a2.annee)

    v.diane = Object.keys(v.diane || {}).reduce((accu, key) => {
        accu = accu.concat(v.diane[key])
        return accu
    }, []).sort((a1, a2) => a1.exercice_diane < a2.exercice_diane)

    emit({siren: this.value.siren, batch: actual_batch}, v) 
  }