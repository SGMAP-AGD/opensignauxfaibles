<template>
  <div>
    <v-card>
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
            <v-bottom-sheet 
              lazy
            >
              <v-btn
                slot="activator"
                flat 
              >
                {{ props.item.effectif }}
              </v-btn>
              <v-card>
                <v-toolbar center card color="indigo lighten-3">
                <h1>EFFECTIFS {{ props.item.raisonSociale }}</h1>
                </v-toolbar>
                <div class="echarts">
                  
                  <IEcharts class="chart"
                  style="height: 500px; width: 1300px"

                  :option="props.item.historyEffectif">
                  </IEcharts>
                </div>
              </v-card>
            </v-bottom-sheet>
          </td>
          <td class="text-xs-right">{{ props.item.dette_urssaf }}</td>
          <td class="text-xs-right">{{ props.item.activite_partielle }}</td>
          <td class="text-xs-right">
            <v-bottom-sheet
              v-if="props.item.all_financiere[props.item.all_financiere.length -1]"
              :close-on-content-click="false"
              offset-y
            >
              <v-btn 
                slot="activator"
                flat
              >
                {{ props.item.all_financiere[props.item.all_financiere.length -1].arrete_bilan }}
              </v-btn>
                    <v-card>
            <v-data-iterator
              :items="props.item.all_financiere"
              hide-actions
              :pagination.sync="pagination"
              content-tag="v-layout"
              row
              wrap
              solid
            >

                  <v-toolbar
                    slot="header"
                    class="mb-2"
                    color="indigo darken-5"
                    dark
                    flat
                  >
                    <v-toolbar-title>Ratios Financiers {{ props.item.raisonSociale }}</v-toolbar-title>
                    <v-spacer></v-spacer>
                    <v-toolbar-items>
                      <v-btn flat @click="fichart = !fichart">
                        <v-icon v-if="fichart">fa-table</v-icon>
                        <v-icon v-if="!fichart">fa-chart-line</v-icon>
                      </v-btn>
                    </v-toolbar-items>
                    </v-toolbar>

                    <v-flex
                      slot="item"
                      slot-scope="fin"
                      xs12
                      sm6
                      md4
                      lg2
                    >
                    <div>
                    <v-card>
                      <v-card-title><h4>{{ fin.item.annee }} ({{ fin.item.arrete_bilan }})</h4></v-card-title>
                      <v-divider></v-divider>
                      <v-list dense>
                        <v-list-tile>
                          <v-list-tile-content>Délai Fournisseur:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.delai_fournisseur }}</v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile>
                          <v-list-tile-content>Dette Fiscale:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.dette_fiscale }}</v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile>
                          <v-list-tile-content>Financier Court Terme:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.financier_court_terme }}</v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile>
                          <v-list-tile-content>Frais Financier:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.frais_financier }}</v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile>
                          <v-list-tile-content>Poids Fond de Roulement:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.poids_frng }}</v-list-tile-content>
                        </v-list-tile>
                        <v-list-tile>
                          <v-list-tile-content>Taux de marge:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ fin.item.taux_marge }}</v-list-tile-content>
                        </v-list-tile>
                      </v-list>
                    </v-card>
                    </div>
                    </v-flex>
                </v-data-iterator>
               </v-card>  
            </v-bottom-sheet>
          </td>
        </tr>
      </template>
    </v-data-table>

    </v-card>
  </div>
</template>

<script>
  import IEcharts from 'vue-echarts-v3/src/full.js'

  export default {
    components: {
      IEcharts
    },
    data: () => ({
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
          text: 'siret',
          align: 'left',
          value: 'siret'
        },
        {
          text: 'raison sociale',
          align: 'left',
          value: 'raison_sociale'
        },
        { text: 'score', value: 'score' },
        { text: 'effectif', value: 'effectif' },
        { text: 'dette urssaf', value: 'dette urssaf' },
        { text: 'activite partielle', value: 'activite_partielle' },
        { text: 'données financières', value: 'données financières' }
      ],
      prediction: []
    }),
    mounted () {
      this.getPrediction()
    },
    methods: {
      getNAF () {
        var self = this
        this.$axios.get(this.$api + '/data/naf').then(response => { self.naf = response.data })
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
        console.log(this.$store.getters.axiosConfig)
        this.$axios.get(this.$api + '/data/prediction/' + this.actualBatch + '/algo1/0', this.$store.getters.axiosConfig).then(response => {
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

<style>
.echarts {
  width: 400px
}
</style>