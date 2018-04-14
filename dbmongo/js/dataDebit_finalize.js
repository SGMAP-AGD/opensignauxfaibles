function finalize(k, v) {
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
                value_dette[debit.periode.start] = (value_dette[debit.periode.start] || []).concat([{"periode": date, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale}])
            })
        }
    })
      
    array_dette = Object.keys(value_dette).reduce((m,p) => {
        var v = {
            "x": value_dette[p].map(function(e) {return e.periode}),
            "y": value_dette[p].map(function(e) {return e.part_ouvriere + e.part_patronale}),
            "name": new Date(p),
            "type": "bar"
        }
        m.push(v)
        return m
    }, [])
      
      
    return array_dette
}