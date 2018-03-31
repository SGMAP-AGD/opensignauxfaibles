function finalize(key, r) {
    // relier les débits
    debit = r.compte.debit
    ecn = {};

    for (i in debit) {
        start = debit[i].periode.start.toISOString().substring(2,7).replace("-","")
        end = debit[i].periode.end.toISOString().substring(2,7).replace("-","")
        num_ecn = debit[i].numero_ecart_negatif
        compte = debit[i].numero_compte
        key = start + "-" + end + "-" + num_ecn + "-" + compte
        if (!(key in ecn)) {
            ecn[key] = []        
        }
        ecn[key].push({
            "hash": i, 
            "numero_historique": debit[i].numero_historique,
            "date_traitement": debit[i].date_traitement
        });
    }

    for (k in ecn) {
        ecn[k].sort(compareDebit)
        for (i = 0; i < ecn[k].length - 1; i++) {
            if (i < ecn[k].length - 1) {
                r.compte.debit[ecn[k][i].hash].debit_suivant = ecn[k][i+1].hash;
            }
        }
    }

    // relier les demandes d'activité partielle aux consommations

    apart = {}
    for (k in r.activite_partielle.demande) {
        var value = r.activite_partielle.demande[k]
        apart[value.id_demande] = {"demande": k,
                                   "consommation": []}
    }

    for (k in r.activite_partielle.consommation) {
        var value = r.activite_partielle.consommation[k]
        if (value.id_conso.substring(0,10) in apart) {
            apart[value.id_conso.substring(0,10)].consommation.push(k)
        }
    }

    for (k in apart) {
        r.activite_partielle.demande[apart[k].demande].hash_consommation = apart[k].consommation
        for (j in apart[k].consommation) {
            r.activite_partielle.consommation[apart[k].consommation[j]].hash_demande = apart[k].demande;
        }
    }
    r.index = {"algo1": Object.keys(r.compte.effectif).some(function() {return true})}
    
    return r;
}