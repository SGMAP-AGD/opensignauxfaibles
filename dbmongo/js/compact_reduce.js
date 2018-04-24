function reduce(key, values) {
    return values.reduce(function (m, o) {
            Object.keys(o.batch).map(function (batch) {
                Object.keys(o.batch[batch]).map(function (t) {
                    m.batch[batch] = (m.batch[batch] || {})
                    m.batch[batch][t] = (m.batch[batch][t] || {})
                    Object.assign(m.batch[batch][t],o.batch[batch][t])
                })
            })
            return m
        })
    } 
