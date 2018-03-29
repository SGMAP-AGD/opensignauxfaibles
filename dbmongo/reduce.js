function reduce(key, values) {
    r = {'compte': {'debit':{}, 'delais': {}, 'cotisation': {}, 'effectif': {}, 'ccsf': {}},
         'altares': {},
         'activite_partielle': {'demande': {}, 'consommation': {}},
         'siret': key}

    for (v in values) {
        r.compte.debit                    = Object.assign(r.compte.debit, values[v].compte.debit)
        r.compte.delais                   = Object.assign(r['compte']['delais'], values[v]['compte']['delais'])
        r.compte.cotisation               = Object.assign(r['compte']['cotisation'], values[v]['compte']['cotisation'])
        r.compte.effectif                 = Object.assign(r['compte']['effectif'], values[v]['compte']['effectif'])
        r.altares                         = Object.assign(r['altares'], values[v]['altares'])
        r.activite_partielle.demande      = Object.assign(r['activite_partielle']['demande'], values[v]['activite_partielle']['demande'])
        r.activite_partielle.consommation = Object.assign(r['activite_partielle']['consommation'], values[v]['activite_partielle']['consommation'])
        }
    
    return r
}