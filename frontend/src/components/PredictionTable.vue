<template>
  <div>
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    :mini-variant = "mini"
    style="z-index: 1"
    >
    <v-list>
      <v-list-tile @click="mini=!mini">
        <v-list-tile-action>
          <v-icon>fa-filter</v-icon>
        </v-list-tile-action>
        <v-list-tile>Sélection</v-list-tile>
      </v-list-tile>
    </v-list>
    </v-navigation-drawer>

    <v-card :class="mini?'widget_large':'widget_tiny'">
    <v-data-table
    v-model="selected"
    :headers="headers"
    :items="prediction"
    :pagination.sync="pagination"
    select-all
    item-key="name"
    class="elevation-1"
    :loading="loading"
    :rows-per-page-items="[10]"
    >
      <template slot="headers" slot-scope="props">
        <tr>
          <th
            v-for="header in props.headers"
            :key="header.text"
            :class="['column sortable', pagination.descending ? 'desc' : 'asc', header.value === pagination.sortBy ? 'active' : '']"
            @click="changeSort(header.value)"
          >
          <v-icon small>arrow_upward</v-icon>
          {{ header.text }}
          </th>
        </tr>
      </template>
      <template slot="items" slot-scope="props">
        <tr :active="props.selected">
          <td>{{ props.item.siret }}</td>
          <td>{{ props.item.raisonSociale }}</td>
          <td class="text-xs-right">{{ Math.round(props.item.score*1000)/1000 }}</td>
          <td class="text-xs-right">
          </td>
          <td class="text-xs-right">{{ (props.item.dette_urssaf>0)?"oui":"non" }}</td>
          <td class="text-xs-right">{{ props.item.activite_partielle }}</td>
        </tr>
      </template>
    </v-data-table>

    </v-card>
  </div>
</template>

<script>
  import IEcharts from 'vue-echarts-v3/src/full.js'

  export default {
    props: ['batchKey'],
    components: {
      IEcharts
    },
    data: () => ({
      mini: true,
      fichart: false,
      loading: true,
      pagination: {
        sortBy: 'name'
      },
      naf: {},
      selected: [],
      actualBatch: '1803',
      headers: [
        {
          text: 'raison sociale',
          align: 'left',
          value: 'raison_sociale'
        },
        { text: 'détection', value: 'score' },
        { text: 'emploi', value: 'emploi' },
        { text: 'urssaf', value: 'urssaf' }
      ],
      prediction: []
    }),
    mounted () {
      this.getPrediction()
    },
    methods: {
      getNAF () {
        var self = this
        this.$axios.get('/api/data/naf').then(response => { self.naf = response.data })
      },
      toggleAll () {
        if (this.selected.length) this.selected = []
        else this.selected = this.desserts.slice()
      },
      changeSort (column) {
        if (this.pagination.sortBy === column) {
          this.pagination.descending = !this.pagination.descending
        } else {
          this.pagination.sortBy = column
          this.pagination.descending = false
        }
      },
      getPrediction () {
        var self = this
        this.loading = true
        this.$axios.get('/api/data/prediction/' + this.batchKey + '/algo1/0', this.$store.getters.axiosConfig).then(response => {
          self.prediction = response.data.map(prediction => {
            var etablissement = self.flattenTypes(
              self.projectBatch(
                prediction.etablissement
                )
            )
            var entreprise = self.flattenTypes(
              self.projectBatch(
                (prediction.entreprise || {})
                )
            )
            var allEffectif = etablissement.effectif.reduce((accu, effectif) => {
              var effectifPeriode = Date.parse(effectif.periode)
              accu[effectifPeriode] = (accu[effectif.periode] || 0) + effectif.effectif
              return accu
            }, {})
            var lastTime = Object.keys(allEffectif).reduce((accu, time) => {
              return (time > accu) ? time : accu
            })
            self.loading = false
            return {
              'siret': prediction._id.siret,
              'score': prediction.score,
              'effectif': allEffectif[lastTime],
              'raisonSociale': (etablissement.sirene || [{'raisonsociale': 'n/a'}])[0].raisonsociale,
              'features': prediction.features,
              'dette_urssaf': prediction.features[prediction.features.length - 1].montant_part_patronale + prediction.features[prediction.features.length - 1].montant_part_ouvriere,
              'activite_partielle': prediction.features[prediction.features.length - 1].apart_heures_consommees,
              'all_financiere': (entreprise.bdf || []).sort((a, b) => a.annee > b.annee).map(bdf => {
                return {
                  'annee': bdf.annee,
                  'secteur': bdf.secteur,
                  'siren': bdf.siren,
                  'arrete_bilan': bdf.arrete_bilan.substring(0, 10),
                  'delai_fournisseur': Math.round(bdf.delai_fournisseur * 1000) / 1000,
                  'dette_fiscale': Math.round(bdf.dette_fiscale * 1000) / 1000,
                  'financier_court_terme': Math.round(bdf.financier_court_terme * 1000) / 1000,
                  'frais_financier': Math.round(bdf.frais_financier * 1000) / 1000,
                  'poids_frng': Math.round(bdf.poids_frng * 1000) / 1000,
                  'taux_marge': Math.round(bdf.taux_marge * 1000) / 1000
                }
              }),
              'historyEffectif': {
                tooltip: {},
                xAxis: [{
                  type: 'time',
                  splitNumber: 3
                }],
                yAxis: {},
                series: [{
                  type: 'line',
                  smooth: true,
                  color: 'blue',
                  data: (Object.keys(allEffectif))
                    .map(key => [key, allEffectif[key]])
                    .sort((a, b) => parseInt(b[0]) - parseInt(a[0]))
                    .map(entry => [new Date(parseInt(entry[0])), entry[1]])
                    .slice(0, 14)
                }]
              }
            }
          })
        })
      },
      projectBatch (o) {
        return Object.keys((o.batch || {})).sort()
          .filter(batch => batch <= this.actualBatch).reduce((m, batch) => {
            Object.keys(o.batch[batch]).forEach((type) => {
              m[type] = (m[type] || {})
              var arrayDelete = (o.batch[batch].compact.delete[type] || [])
              if (arrayDelete !== {}) {
                arrayDelete.forEach(hash => {
                  delete m[type][hash]
                })
              }
              Object.assign(m[type], o.batch[batch][type])
            })
            return m
          }, {})
      },
      flattenTypes (o) {
        return Object.keys(o).filter(type => type !== 'compact').reduce((accu, type) => {
          accu[type] = Object.values(o[type])
          return accu
        }, {})
      }
    }
  }
</script>

<style scoped>
.echarts {
  width: 400px
}
.widget_tiny {
  position: absolute;
  left: 320px;
  top: 20px; 
  right: 20px;
}
.widget_large {
  position: absolute;
  left: 100px;
  top: 20px;
  right: 20px;
}
</style>