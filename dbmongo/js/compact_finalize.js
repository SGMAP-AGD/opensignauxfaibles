function finalize(k, o) {
    var deleteOld = new Set(["effectif", "apdemande", "apconso"])
    Object.keys((o.batch||{})).sort().reduce((m, batch) => {
        Object.keys(o.batch[batch]).map(type => {
            m[type] = (m[type] || new Set())
            var keys = Object.keys(o.batch[batch][type])
            if (deleteOld.has(type) && o.batch[batch].compact.status == false) {
                var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
                discardKeys.forEach(key => {
                    m[type].delete(key)
                    o.batch[batch][type][key] = null
                })
            }
            keys.filter(key => (m[type].has(key))).forEach(key => delete o.batch[batch][type][key])
            m[type] = new Set([...m[type]].concat(keys))
            o.batch[batch].compact.status = true
        })
        return m
    }, {})
    Object.keys((o.etablissement||{})).forEach(etablissement => {
        Object.keys(o.etablissement[etablissement].batch).sort().reduce((m, batch) => {
            Object.keys(o.etablissement[etablissement].batch[batch]).map(type => {
                m[type] = (m[type] || new Set())
                var keys = Object.keys(o.etablissement[etablissement].batch[batch][type])
                if (deleteOld.has(type) && o.etablissement[etablissement].batch[batch].compact.status == false) {
                    var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
                    discardKeys.forEach(key => {
                        m[type].delete(key)
                        o.etablissement[etablissement].batch[batch][type][key] = null
                    })
                }
                keys.filter(key => (m[type].has(key))).forEach(key => delete o.etablissement[etablissement].batch[batch][type][key])
                m[type] = new Set([...m[type]].concat(keys))
                o.etablissement[etablissement].batch[batch].compact.status = true
            })
            return m
        }, {})
    })

    o.index = {"algo1": false}
    Object.keys(o.etablissement).forEach(etablissement => {
        Object.keys(o.etablissement[etablissement].batch).forEach(batch => {
            Object.keys((o.etablissement[etablissement].batch[batch].effectif||{})).forEach(effectif => o.index.algo1 = true)      
        })
    })

    o.region = []
    Object.keys(o.etablissement).forEach(etablissement => {
        Object.keys(o.etablissement[etablissement].batch).forEach(batch => {
            Object.keys((o.etablissement[etablissement].batch[batch].sirene||{})).forEach(hash => {
                r = o.etablissement[etablissement].batch[batch].sirene[hash].region
                if (o.region.indexOf(r) == -1 && r != "") {
                    o.region.push(r)
                }
            })
        })
    })

    return o
}