function map() { 
    var m = Object.keys((this.value.batch || {})).sort().reduce((accu, batch) => {
        Object.keys(this.value.batch[batch]).forEach((type) => {
            accu[type] = (accu[type] || {})
            Object.assign(accu[type], this.value.batch[batch][type])
        })
        return accu
    }, { "siren": this.value.siret.substring(0,9) })

    emit(this.value.siret, m) 
}