function finalize(k, v) {

    v = Object.keys((v.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
        Object.keys(v.batch[batch]).forEach((type) => {
            m[type] = (m[type] || {})
            var  array_delete = (v.batch[batch].compact.delete[type]||[])
            if (array_delete != {}) {array_delete.forEach(hash => {
                delete m[type][hash]
            })
            }
            Object.assign(m[type], v.batch[batch][type])
        })
        return m
    }, { "siren": k })

    v = Object.keys(v).filter(k => k !== 'siren' && k !== 'compact').reduce((accu, type) => {
        accu[type] = Object.keys(v[type]).map(k => v[type][k])
        return accu
    }, {})

    return v
}