function finalize(k, o) {
    var deleteOld = new Set(["effectif", "apdemande", "apconso"])
    Object.keys(o.batch).sort().reduce((m, batch) => {
        Object.keys(o.batch[batch]).map(type => {
            m[type] = (m[type] || new Set())
            var keys = Object.keys(o.batch[batch][type])
            if (deleteOld.has(type) && o.batch[batch].compact.status == false) {
                var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
                discardKeys.forEach(key => {
                    m[type].delete(key)
                    o.batch[batch][type][key] = null
                })
            }
            keys.filter(key => (m[type].has(key))).forEach(key => delete o.batch[batch][type][key])
            m[type] = new Set([...m[type]].concat(keys))
            o.batch[batch].compact.status = true
        })
        return m
    }, {})


        // // relier les débits
        // var debit = Object.keys(o.batch).reduce((m, batch) => {
        //     Object.assign(m, (o.batch[batch].debit || {}))
        //     return m
        //     },{})

        // var ecn = Object.keys(debit).reduce((m, h) => {
        //     var d = [h, debit[h]]
        //     var start = d[1].periode.start
        //     var end = d[1].periode.end
        //     var num_ecn = d[1].numero_ecart_negatif
        //     var compte = d[1].numero_compte
        //     var key = start + "-" + end + "-" + num_ecn + "-" + compte
        //     m[key] = (m[key] || []).concat([{
        //         "hash": d[0], 
        //         "numero_historique": d[1].numero_historique,
        //         "date_traitement": d[1].date_traitement
        //     }])
        //     return m 
        // }, {})

        // Object.keys(ecn).forEach(i => {
        //     ecn[i].sort(compareDebit)
        //     var l = ecn[i].length
        //     ecn[i].forEach((e,idx) => {
        //         if (idx < l-2) {
        //             console.log(idx)
        //             debit[e.hash].debit_suivant = ecn[idx+1].hash;  
        //         }
        //     })
        // })

    // relier les demandes d'activité partielle aux consommations

    o.index = {"algo1": Object.keys(o.batch).some(batch => o.batch[batch].effectif)}
    
    return o
}