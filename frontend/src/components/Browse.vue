<template>
<div>
  <v-toolbar height="35px" class="toolbar" color="#ffffff"  app>
    <v-icon
     @click="drawer=!drawer"
    class="fa-rotate-180"
    v-if="!drawer"
    color="primary"
    key="toolbar"
    >mdi-backburger</v-icon>
    <div style="width: 100%; text-align: center;"  class="titre">
      Détection
    </div>
    <v-spacer></v-spacer>
    <v-icon color="primary" @click="rightDrawer=!rightDrawer">mdi-magnify</v-icon>
  </v-toolbar>
  <div style="width:100%">
  <v-navigation-drawer :class="(rightDrawer?'elevation-6':'') + 'rightDrawer'" v-model="rightDrawer" right app>
    <v-toolbar flat class="transparent">
      <v-list class="pa-0">
        <v-list-tile avatar>
          <v-list-tile-avatar>
            <img src="/static/logo_signaux_faibles_cercle.svg">
          </v-list-tile-avatar>

          <v-list-tile-content>
            <v-list-tile-title><span class="fblue">Signaux</span>·<span class="fred">Faibles</span></v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-avatar>
            <v-icon @click="drawer=!drawer">mdi-chevron-left</v-icon>
          </v-list-tile-avatar>
        </v-list-tile>
      </v-list>
    </v-toolbar>
     <v-list class="pt-0" dense>
        <v-divider></v-divider>
        <v-list-tile to="/">
          <v-list-tile-action>
            <v-icon>home</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Accueil</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile to="/Browse">
          <v-list-tile-action>
            <v-icon>fa-search</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Détection</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile to="/data">
          <v-list-tile-action>
            <v-icon>fa-database</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Données</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

        <v-list-tile to="/admin">
          <v-list-tile-action>
            <v-icon>fa-cog</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Administration</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-divider></v-divider>
        <v-list-tile @click="logout()">
          <v-list-tile-action>
            <v-icon>logout</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title>Se déconnecter</v-list-tile-title>
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
</div>
</template>

<script>
export default {
  data () {
    return {
      active: 0
    }
  },
  mounted () {
    this.$store.commit('updateBatches')
  },
  methods: {
    setCurrentBatchKey (batchKey) {
      this.currentBatchKey = batchKey
    },
    close (tabIndex) {
      this.activeTab = Math.min(this.activeTab, (this.tabs.length - 2))
      this.tabs = this.tabs.filter((tab, index) => index !== tabIndex)
    }
  },
  computed: {
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
    tabs: {
      get () { return this.$store.getters.getTabs },
      set (tabs) { this.$store.dispatch('updateTabs', tabs) }
    },
    activeTab: {
      get () { return this.$store.getters.activeTab },
      set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
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
    }
  },
  name: 'Browse'
}
</script>

<style>
div.titre {
  color: #20459a;
  font-family: 'Signika', sans-serif;
  font-weight: 500;
  color: primary;
  font-size: 18px;
}
</style>
