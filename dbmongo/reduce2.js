function reduce(key, values) {
    r = {"altares":    Math.min(1, Object.keys(this.value.altares).length),
         "effectif":   Math.min(1, Object.keys(this.value.compte.effectif).length),
         "debits":     Math.min(1, Object.keys(this.value.compte.debit).length),
         "cotisation": Math.min(1, Object.keys(this.value.compte.cotisation).length),
         "delais":     Math.min(1, Object.keys(this.value.compte.delais).length)};
    return r
}