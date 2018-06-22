function reduce(key, values) {
    return values.reduce((m, v) => {
        Object.assign(m, v)
        return m
    }, {})
}