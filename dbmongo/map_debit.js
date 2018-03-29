function map() {
    a = [];
    debit = this.value.compte.debit;
    ecn = {};
    print(Object.keys(debit).length)
    for (i in debit) {
        start = debit[i]["periode"]["start"].toISOString().substring(2,7).replace("-","")
        end = debit[i]["periode"]["end"].toISOString().substring(2,7).replace("-","")
        num_ecn = debit[i]["numero_ecart_negatif"]
        compte = debit[i]["numero_compte"]
        key = start + "-" + end + "-" + num_ecn + "-" + compte
        if (!(key in ecn)) {
            ecn[key] = []        
        }
        ecn[key].push(debit[i]);
    }

    for (k in ecn) {
            ecn[k].sort(compareDebit)
    }

    emit(this.value.siret,
        ecn)
}