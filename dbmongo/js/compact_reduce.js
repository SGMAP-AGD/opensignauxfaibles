function reduce(key, values) {
    return values.reduce((m, value) => {
        Object.keys((value.batch||{})).forEach(batch => {
            Object.keys(value.batch[batch]).forEach(type => {
                m.batch = (m.batch||{})
                m.batch[batch] = (m.batch[batch] || {})
                m.batch[batch][type] = (m.batch[batch][type] || {})
                Object.assign(m.batch[batch][type],value.batch[batch][type])
            })
        })

        Object.keys((value.etablissement||{})).forEach(etablissement => {
            Object.keys(value.etablissement[etablissement].batch).map(batch => {
                Object.keys(value.etablissement[etablissement].batch[batch]).map(type => {
                    m.etablissement = (m.etablissement||{})
                    m.etablissement[etablissement] = (m.etablissement[etablissement]||{})
                    m.etablissement[etablissement].batch = (m.etablissement[etablissement].batch||{})
                    m.etablissement[etablissement].batch[batch] = (m.etablissement[etablissement].batch[batch] || {})
                    m.etablissement[etablissement].batch[batch][type] = (m.etablissement[etablissement].batch[batch][type] || {})
                    Object.assign(m.etablissement[etablissement].batch[batch][type],value.etablissement[etablissement].batch[batch][type])
                })
            })
        })


        return m
    },{"siren": key})

}