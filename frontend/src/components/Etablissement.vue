<template>
  <div>
    <div class="widget">
      <v-card>
        <v-toolbar class="headline toolbar elevation-3" color='indigo lighten-4' card>
          Etablissement {{ sirene.raisonsociale }}
          <v-spacer/>
          <v-icon 
          color="red"
          @click="close()">fa-times-circle</v-icon>
        </v-toolbar>
        <v-card-title>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-layout wrap>
              <v-flex xs6>
                {{ sirene.numvoie }} RUE {{ sirene.typevoie  }}<br/>
                {{ sirene.codepostal }} {{ sirene.commune }}<br/>
                <br/>
                siret: {{ sirene.siret }} {{ siret }} <br/>
                {{ sirene.naturejuridique }}<br/>
                Création: {{ printDate(sirene.debut_activite) }}<br/>
                <br/>
                {{ naf.n1[naf.n5to1[sirene.ape]] }}<br/>
                {{ (naf.n5 || {})[sirene.ape] }}<br/>
                Code APE: {{ sirene.ape }}<br/>
              </v-flex>
              <v-flex xs6>
                <iframe width="400" height="300" frameborder="0" scrolling="no" marginheight="0" marginwidth="0" :src="'https://www.openstreetmap.org/export/embed.html?bbox=' + (sirene.longitude - 0.05) + '%2C' + (sirene.lattitude  - 0.05) + '%2C' + (sirene.longitude + 0.05) + '%2C' + (sirene.lattitude + 0.05) + '&amp;layer=mapnik&amp;marker=' + sirene.lattitude + '%2C' + sirene.longitude" style="border: 1px solid black"></iframe><br/><small><a href="https://www.openstreetmap.org/#map=19/47.31581/5.05088">Afficher une carte plus grande</a></small>
              </v-flex>
              <v-flex xs12>
                <v-data-iterator
                  :items="zipDianeBDF"
                  :rows-per-page-items="[4]"
                  :pagination.sync="pagination"
                  content-tag="v-layout"
                  row
                  wrap>
                  <v-flex
                    slot="item"
                    slot-scope="props"
                    xs12
                    sm6
                    md4
                    lg3
                  >
                    <v-card>
                      <v-card-title class="subheading font-weight-bold">{{ props.item.annee }}</v-card-title>

                      <v-divider></v-divider>

                      <v-list dense>
                        <v-list-tile>
                          <v-list-tile-content>Arrété Bilan:</v-list-tile-content>
                          <v-list-tile-content class="align-end">{{ (props.item.bdf[0]||{})['arrete_bilan'].substring(0,10) }}</v-list-tile-content>
                        </v-list-tile>
                      </v-list>
                    </v-card>
                  </v-flex>
                </v-data-iterator>
              </v-flex>
            </v-layout>
          </v-container>
          {{ zipDianeBDF }}
        </v-card-text>
      </v-card>
      
    </div>
  </div>
</template>

<script>
export default {
  props: ['siret', 'batch'],
  name: 'Etablissement',
  data () {
    return {
      bilan: true,
      urssaf: true,
      apart: true,
      etablissement: {},
      entreprise: {},
      pagination: null
    }
  },
  methods: {
    close () {
      this.tabs = this.tabs.filter((tab, index) => index !== this.activeTab)
      this.activeTab = this.activeTab - 1
    },
    printDate (date) {
      return date.substring(0, 10)
    }
  },
  mounted () {
    this.$axios.get('/api/data/etablissement/1802/' + this.siret).then(response => {
      this.etablissement = response.data.etablissement[0].value
      this.entreprise = response.data.entreprise[0].value
    })
    this.$axios.get('/api/data/naf').then(response => { this.naf = response.data })
  },
  computed: {
    activeTab: {
      get () { return this.$store.getters.activeTab },
      set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
    },
    tabs: {
      get () { return this.$store.getters.getTabs },
      set (tabs) { this.$store.dispatch('updateTabs', tabs) }
    },
    sirene () {
      return ((this.etablissement.sirene || [])[0]) || {}
    },
    bdf () {
      return (this.entreprise.bdf || []).sort((a, b) => a.annee > b.annee).reverse()
    },
    diane () {
      return (Object.keys(this.entreprise.diane || []).map(k => this.entreprise.diane[k]) || [null]).sort((a, b) => a.annee > b.annee).reverse()
    },
    zipDianeBDF () {
      let annees = new Set(this.bdf.map(b => b.annee).concat(this.diane.map(d => d.annee)))
      return Array.from(annees).sort((a, b) => a < b).map(a => {
        return {
          annee: a,
          bdf: this.bdf.filter(b => b.annee === a),
          diane: this.diane.filter(d => d.annee === a)
        }
      })
    }
  }
}
</script>

<style scoped>
.echarts {
  width: 400px
}
.widget {
  position: absolute;
  left: 20px;
  top: 20px; 
  right: 20px;
}
</style>