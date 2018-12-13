function map() { 
  v = Object.keys((this.value.batch || {}))
  .sort()
  .reduce((m, batch) => {
    Object.keys(this.value.batch[batch])
    .filter(type => ['sirene'].includes(type))
    .forEach((type) => {
      m[type] = (m[type] || {})
      var  array_delete = (this.value.batch[batch].compact.delete[type]||[])
      if (array_delete != {}) {array_delete.forEach(hash => {
        delete m[type][hash]
      })
      }
      Object.assign(m[type], this.value.batch[batch][type])
    })
    return m
  }, {})

  sirene = Object.keys(v.sirene || {}).reduce((accu, k) => {
    accu.raisonsociale = v.sirene[k].raisonsociale
    accu.nicsiege = v.sirene[k].nicsiege
    accu.adresse = v.sirene[k].adresse
    accu.gps = [v.sirene[k].lattitude, v.sirene[k].longitude]
    accu.creation = v.sirene[k].creation
    accu.ape = v.sirene[k].ape
    return accu
  }, {})

  r = {
    sirene
  }
  emit(this.value.siret, r) 
}