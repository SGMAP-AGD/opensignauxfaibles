function map() { 
    v = Object.keys((this.value.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
        Object.keys(this.value.batch[batch]).forEach((type) => {
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

    emit({siren: this.value.siren, batch: actual_batch}, v) 
  }