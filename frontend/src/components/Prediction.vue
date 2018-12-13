<template>
<div>
  <v-toolbar
    height="35px"
    class="toolbar elevation-12"
    color="#ffffff"
    extension-height="48px"
    app
  >
    <v-icon
      @click="drawer=!drawer"
      class="fa-rotate-180"
      v-if="!drawer"
      color="#e0e0ffef"
      key="toolbar"
    >
      mdi-backburger
    </v-icon>
    <div style="width: 100%; text-align: center;"  class="titre">
      Détection
    </div>
    <v-spacer></v-spacer>
    <v-icon color="#e0e0ffef" v-if="!rightDrawer" @click="rightDrawer=!rightDrawer">mdi-target</v-icon>
  </v-toolbar>
  
  <div style="width:100%">
    <v-navigation-drawer :class="(rightDrawer?'elevation-6':'') + 'rightDrawer'" v-model="rightDrawer" right app>
      <v-toolbar flat class="transparent" height="100">
        <v-list class="pa-0">
          <v-list-tile avatar>
            <v-list-tile-avatar>
              <v-icon @click="rightDrawer=!rightDrawer">mdi-target</v-icon>
          </v-list-tile-avatar>
          <v-spacer></v-spacer>

          
            <v-img src="/static/regions/PDL.svg"></v-img>
          
          </v-list-tile>
        </v-list>
      </v-toolbar>
      <v-list two-line>
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
      </v-list>
      <v-footer class="elevation-6" style="color: blue; width:100%; position: fixed; bottom: 0px;">
        <v-btn
          flat
          icon
          color="blue"
          href="https://github.com/entrepreneur-interet-general/opensignauxfaibles">
          <v-icon>fab fa-github</v-icon>
        </v-btn>
      </v-footer>
    </v-navigation-drawer>
  </div>
  <PredictionWidget v-for="p in prediction" :key="p._id.siret" :prediction="p"/>
 
</div>
</template>

<script>
import PredictionWidget from '@/components/PredictionWidget'
export default {
  data () {
    return {
      effectifClass: [10, 20, 50, 100],
      active: 0,
      prediction: [],
      naf: 'Industrie manufacturière',
      minEffectif: 20,
      entrepriseConnue: true,
      horsCCSF: true,
      horsProcol: true,
      loading: false,
    }
  },
  mounted () {
    this.getPrediction()
    this.$store.dispatch('getNAF')
    this.$store.commit('updateBatches')
  },
  methods: {
    getPrediction() {
      this.loading = true;
      var self = this;
      this.$axios.post('/api/data/prediction').then(response => {
        var prediction = response.data;
        prediction.forEach(p => {
          p.bdf = Object.keys(p.bdf || {})
            .map(b => p.bdf[b])
            .sort((a, b) => a.annee < b.annee);
        });
        this.prediction = prediction;
        self.loading = false;
      });
    }
  },
  computed: {
    naf1 () {
      return Object.keys(this.$store.state.naf.n1 || {}).sort().map(n => n + ' - ' + this.$store.state.naf.n1[n].substring(0,60) )
    },
    scrollTop () {
      return this.$store.state.scrollTop
    },
    height: {
      get () {
        return this.$store.state.height
      },
      set (height) {
        this.$store.dispatch('setHeight', height)
      }
    },
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    },
    rightDrawer: {
      get () {
        return this.$store.state.rightDrawer
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
    currentBatchKey: {
      get () {
        return this.$store.state.currentBatchKey
      },
      set (value) {
        this.$store.commit('setCurrentBatchKey', value)
      }
    },
    batches () {
      return this.$store.state.batches.filter(b => b.readonly === true).map(batch => batch.id.key)
    },
    visiblePrediction () {

    },
    detectionLength () {
      return Math.round((this.height + this.scrollTop) / 1000 + 1)*10 ;
    }
  },
  components: { PredictionWidget },
  name: 'Browse'
}
</script>

<style>
div.titre {
  color: #e0e0ffef;
  font-family: 'Signika', sans-serif;
  font-weight: 500;
  color: primary;
  font-size: 18px;
}
</style>
