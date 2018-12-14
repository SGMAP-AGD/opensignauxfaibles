function map() { 
  v = Object.keys((this.value.batch || {}))
  .filter(b => b <= actual_batch)
  .sort()
  .reduce((m, batch) => {
    Object.keys(this.value.batch[batch])
    .filter(type => ['sirene', 'effectif'].includes(type))
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

  var sirene = Object.keys(v.sirene || {}).reduce((accu, k) => {
    accu.raisonsociale = v.sirene[k].raisonsociale
    accu.nicsiege = v.sirene[k].nicsiege
    accu.adresse = v.sirene[k].adresse
    accu.gps = [v.sirene[k].lattitude, v.sirene[k].longitude]
    accu.creation = v.sirene[k].creation
    accu.ape = v.sirene[k].ape
    return accu
  }, {})

  var effectif = Object.keys(v.effectif || {}).reduce((accu, k) => {
    if (v.effectif[k].periode.getTime() > accu.date.getTime()) {
      accu.date = v.effectif[k].periode
      accu.effectif = v.effectif[k].effectif
    }
    return accu
  }, {date: new Date(0)})

  r = {
    sirene,
    effectif
  }
  emit({siret: this.value.siret, batch: actual_batch}, r) 
}