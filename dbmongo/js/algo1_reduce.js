function reduce(key, values) {
    return values.reduce((m, value) => {
        Object.keys(value.batch).sort().map(batch => {
            Object.keys(value.batch[batch]).map(type => {
                m.batch[batch] = (m.batch[batch] || {})
                m.batch[batch][type] = (m.batch[batch][type] || {})
                Object.assign(m.batch[batch][type],value.batch[batch][type])
            })
        })
        return m
    })
}