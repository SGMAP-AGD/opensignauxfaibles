<template>
<div>
    <md-card>
      <md-card-header>
        <div class="md-title">Téléchargement des données entreprises</div>

      </md-card-header>
      <md-card-content>

          <label>Insérer les siret souhaités:</label>
          <md-chips v-model="sirets" md-placeholder="Ajouter siret..."></md-chips>

        <md-button class='md-raised md-primary' v-on:click='getCotisation()'>Télécharger</md-button>

      </md-card-content>
    </md-card>
</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'DataView',
  data () {
    return {
      sirets: [],
      data: {},
      effectif: [],
      cotisation: [],
      debit: [],
      textarea: ''
    }
  },
  methods: {
    getCotisation () {
      var self = this
      axios.post(this.$api + '/data', this.sirets).then(function (response) {
        self.data = response.data

        self.data.forEach(dt => {
          var tmp = Object.keys(dt.value.map_effectif).map(key => {
            return { 'siret': dt.value.siret, 'periode': key, 'effectif': dt.value.map_effectif[key] }
          })
          tmp.forEach(t => {
            self.effectif.push(t)
          })

          Object.keys(dt.value.value_cotisation).forEach(key => {
            dt.value.value_cotisation[key].forEach(c => {
              c.siret = dt.value.siret
              self.cotisation.push(c)
            })
          })

          Object.keys(dt.value.value_dette).forEach(d => {
            dt.value.value_dette[d].forEach(e => {
              e.periode_algorithme = d.substring(0, 10)
              e.siret = dt.value.siret
              self.debit.push(e)
            })
          })
        })

        let contentDebit = 'siret;periode_algorithme;periode_debit;part_ouvriere;part_patronale\n'
        self.debit.forEach(d => {
          let row = d.siret + ',' +
            d.periode_algorithme.substring(0, 10) + ',' +
            d.periode.substring(0, 10) + ',' +
            d.part_ouvriere + ',' +
            d.part_patronale
          contentDebit += row + '\n'
        })

        let contentCotisation = 'siret;periode_algorithme;periode_debit;periode_debut_cotisation;periode_fin_cotisation;numero_compte;montant_cotisation;duree_periode;ecriture\n'
        self.cotisation.forEach(c => {
          let row = c.siret + ',' +
          c.periode_algorithme.substring(0, 10) + ',' +
          c.periode_debit + ',' +
          c.periode_debut_cotisation.substring(0, 10) + ',' +
          c.periode_fin_cotisation.substring(0, 10) + ',' +
          c.numero_compte + ',' +
          c.montant_cotisation + ',' +
          c.duree_periode + ',' +
          c.ecriture
          contentCotisation += row + '\n'
        })
        let contentEffectif = 'siret;periode_algorithme;effectif\n'
        self.effectif.forEach(e => {
          let row = e.siret + ',' +
          e.periode.substring(0, 10) + ',' +
          e.effectif
          contentEffectif += row + '\n'
        })
        self.download('debit.csv', contentDebit)
        self.download('cotisation.csv', contentCotisation)
        self.download('effectif.csv', contentEffectif)
      })
    },
    download (filename, text) {
      var element = document.createElement('a')
      element.setAttribute('href', 'data:text/csv;charset=utf-8,' + encodeURIComponent(text))
      element.setAttribute('download', filename)
      element.style.display = 'none'
      document.body.appendChild(element)
      element.click()
      document.body.removeChild(element)
    },
    isSiret (str) {
      return str.length === 14 && !(isNaN(str))
    }
  }
}

</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
