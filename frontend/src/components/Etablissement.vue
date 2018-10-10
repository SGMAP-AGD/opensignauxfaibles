<template>
  <div>
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    style="z-index: 1"
    >
      <v-list>
        <v-toolbar class="elevation-1">
          Etablissement {{ siret }}
          <v-spacer/>
        </v-toolbar>
        <v-list-tile>
          <v-list-tile-action>
            <v-checkbox v-model="bilan"/>
          </v-list-tile-action>
          <v-list-tile-content>
            Informations financières
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-checkbox v-model="urssaf"/>
          </v-list-tile-action>
          <v-list-tile-content>
            Cotisations sociales
          </v-list-tile-content> 
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-checkbox v-model="urssaf"/>
          </v-list-tile-action>
          <v-list-tile-content>
            Activité partielle
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    
    
    </v-navigation-drawer>
    <div class="widget">
      <v-card>
        <v-toolbar card>
          Etablissement {{ siret }}
          <v-spacer/>
          <v-icon 
          color="red"
          @click="close()">fa-times-circle</v-icon>
        </v-toolbar>
        <v-card-title>
          {{ JSON.stringify(etablissement, null, 2)}}
        </v-card-title>
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
      entreprise: {}
    }
  },
  methods: {
    close () {
      this.tabs = this.tabs.filter((tab, index) => index !== this.activeTab)
      this.activeTab = this.activeTab - 1
    }
  },
  computed: {
    activeTab: {
      get () { return this.$store.getters.activeTab },
      set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
    },
    tabs: {
      get () { return this.$store.getters.getTabs },
      set (tabs) { this.$store.dispatch('updateTabs', tabs) }
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
  left: 320px;
  top: 20px; 
  right: 20px;
}
</style>