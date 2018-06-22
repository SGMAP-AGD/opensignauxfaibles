function finalize(k, v) {
    array_periode = serie_periode.reduce((accu,periode) => {
        accu[periode.getTime()] = {"periode": periode}
        return accu
    }, {})
    
    

    return_value = { "siren": k.substring(0, 9)}
    return_value[k] = Object.keys(array_periode).sort().map(t => array_periode[t])
    return return_value
}