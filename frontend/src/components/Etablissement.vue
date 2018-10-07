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
          Etablissement {{ param }}
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
          Etablissement {{ param }}
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
      etablissement: {}
    }
  },
  mounted () {
    this.$axios.get('/api/data/etablissement/' + this.batch + '/' + this.siret).then(response => {
      this.etablissement = response.data
    })
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