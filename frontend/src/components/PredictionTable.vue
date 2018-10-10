<template>
  <div>
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    :mini-variant = "mini"
    style="z-index: 1"
    >
    <v-list two-line>
      <v-list-group
      v-model="nomini">
        <v-list-tile  slot="activator" @click="mini=!mini">
          <v-list-tile-action>
            <v-icon>fa-filter</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>Filtrage</v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-industry</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="naf1"
              v-model="naf"
              label="Secteur d'activité"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-users</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="effectifClass"
              v-model="minEffectif"
              label="Effectif minimum"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="entrepriseConnue">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            Entreprise non suivie
          </v-list-tile-content>

        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="horsCCSF">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors CCSF
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="horsProcol">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors Procédure Collective
          </v-list-tile-content>
        </v-list-tile>
        </v-list-group>
    </v-list>
    </v-navigation-drawer>

    <v-card :class="mini?'widget_large':'widget_tiny'">
    <v-data-table
    v-model="selected"
    :headers="headers"
    :items="predictionFiltered"
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
        <tr 
        :active="props.selected"
        @click="open(props.item)"
        >
          
          <td>{{ props.item.raison_sociale }} </td>
          <td class="text-xs-center"><widgetPrediction :prob="props.item.prob" :diff="props.item.diff"/></td>
          <td class="text-xs-center">{{ props.item.departement }}</td>
          <td class="text-xs-right">
            {{ props.item.effectif }}
          </td>
          <td class="text-xs-right">{{ props.item.default_urssaf?"oui":"non" }}</td>

        </tr>
      </template>
    </v-data-table>

    </v-card>
  </div>
</template>

<script>
  import IEcharts from 'vue-echarts-v3/src/full.js'
  import widgetPrediction from '@/components/widgetPrediction'
  export default {
    props: ['batchKey'],
    components: {
      IEcharts,
      widgetPrediction
    },
    data: () => ({
      effectifClass: [10, 20, 50, 100],
      selected: [],
      mini: true,
      loading: true,
      pagination: {
        sortBy: 'name'
      },
      naf1: [
        'Tous',
        'Activités spécialisées, scientifiques et techniques',
        'Activités de services administratifs et de soutien',
        'Industrie manufacturière',
        'Hébergement et restauration',
        'Construction',
        'Transports et entreposage',
        'Commerce ; réparation d\'automobiles et de motocycles',
        'Santé humaine et action sociale',
        'Autres activités de services',
        'Arts, spectacles et activités récréatives',
        'Industries extractives',
        'Production et distribution d\'eau ; assainissement, gestion des déchets et dépollution',
        'Information et communication',
        'Activités financières et d\'assurance',
        'Activités immobilières',
        'Agriculture, sylviculture et pêche',
        'Production et distribution d\'électricité, de gaz, de vapeur et d\'air conditionné',
        'Activités extra-territoriales'
      ],
      headers: [
        {
          text: 'raison sociale',
          align: 'left',
          value: 'raison_sociale'
        },
        { text: 'détection', value: 'prob' },
        { text: 'département', value: 'departement' },
        { text: 'emploi', value: 'effectif' },
        { text: 'Défault urssaf', value: 'default_urssaf' }
      ],
      prediction: [],
      naf: 'Industrie manufacturière',
      minEffectif: 20,
      entrepriseConnue: true,
      horsCCSF: true,
      horsProcol: true
    }),
    computed: {
      nomini: {
        get () { return !this.mini },
        set (mini) { this.mini = !mini }
      },
      predictionFiltered () {
        return this.prediction.filter(p => this.applyFilter(p))
      },
      tabs: {
        get () { return this.$store.getters.getTabs },
        set (tabs) { this.$store.dispatch('updateTabs', tabs) }
      }
    },
    mounted () {
      this.getPrediction()
    },
    methods: {
      open (etab) {
        if (this.tabs.findIndex(t => t.siret === etab._id.siret) === -1) {
          this.tabs.push({
            'type': 'Etablissement',
            'param': etab.raison_sociale,
            'siret': etab._id.siret,
            'batch': '1802'
          })
        }
      },
      applyFilter (p) {
        return (this.naf === 'Tous' || p.naf1 === this.naf) &&
        (p.effectif >= this.minEffectif) &&
        (p.connu === false || this.entrepriseConnue === false) &&
        (p.ccsf === false || this.horsCCSF === false) &&
        (p.procol === 'in_bonis' || this.horsProcol === false)
      },
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
        this.loading = true
        var self = this
        this.$axios.get('/api/data/prediction').then(response => {
          this.prediction = response.data
          self.loading = false
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