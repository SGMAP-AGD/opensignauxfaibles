function finalize(k, o) {
    var deleteOld = new Set(["effectif", "apdemande", "apconso", "ccsf"])

    Object.keys((o.batch||{})).sort().reduce((m, batch) => {
        
        Object.keys(o.batch[batch]).map(type => {
            m[type] = (m[type] || new Set())
            var keys = Object.keys(o.batch[batch][type])
            if (deleteOld.has(type) && (o.batch[batch].compact||{"status": false}) == false) {
                var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
                discardKeys.forEach(key => {
                    m[type].delete(key)
                    delete o.batch[batch][type][key] 
                    o.batch[batch].compact.delete = (o.batch[batch].compact.delete||{})
                    o.batch[batch].compact.delete[type] = (o.batch[batch].compact.delete[type]||[]).concat([key])
                })
            }

            keys.filter(key => (m[type].has(key))).forEach(key => delete o.batch[batch][type][key])
            m[type] = new Set([...m[type]].concat(keys))

        })

        o.batch[batch].compact = (o.batch[batch].compact||{})
        o.batch[batch].compact.status = true

        return m
        
    }, {})

    o.index = {"algo1": false}
    Object.keys(o.batch).forEach(batch => {
        Object.keys((o.batch[batch].effectif||{})).forEach(effectif => o.index.algo1 = true)      
    })

    return o
}