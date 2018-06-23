function map() { 
    var value = Object.keys((this.value.batch || {})).sort().reduce((m, batch) => {
        Object.keys(this.value.batch[batch]).forEach((type) => {
            m[type] = (m[type] || {})
            Object.assign(m[type], this.value.batch[batch][type])
        })
        return m
    }, { "siren": this.value.siren })

    emit(this.value.siren, value) 
}