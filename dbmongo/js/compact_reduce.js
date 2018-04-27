function reduce(key, values) {
    r = values.reduce((m, value) => {
        Object.keys(value.batch).map(batch => {
            Object.keys(value.batch[batch]).map(type => {
                m.batch[batch] = (m.batch[batch] || {})
                m.batch[batch][type] = (m.batch[batch][type] || {})
                Object.assign(m.batch[batch][type],value.batch[batch][type])
            })
        })
        return m
    })
    return r
}