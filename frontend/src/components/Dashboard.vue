<template>
<div>
  <v-toolbar class="toolbar" color="#c9aec5" height="35px" app>
    <v-icon
    @click="drawer=!drawer"
    class="fa-rotate-180"
    v-if="!drawer"
    color="#ffffff"
    key="toolbar"
    >mdi-backburger
    </v-icon>
    <div style="width: 100%; text-align: center;" class="titre">
      Tableau de Bord
    </div>
    <v-spacer></v-spacer>
    <v-icon color="#ffffff" @click="rightDrawer=!rightDrawer">mdi-view-dashboard</v-icon>
  </v-toolbar>
  <v-container fluid grid-list-md text-xs-center>
      <v-layout row wrap justify-start>
        <v-flex xs12></v-flex>
        <v-flex xs6>
          <v-card>
            <v-card-title class="headline font-weight-thin" >
              Calcul prédictions
            </v-card-title>
            <v-divider/>
            <v-card-text>
              <v-list>
                <v-list-tile avatar>
                  Effectué: x%
                </v-list-tile>
                <v-list-tile avatar>
                  En attente: y%
                </v-list-tile>
              </v-list>
            </v-card-text>
          </v-card>
        </v-flex>
        <v-flex xs6>
          <v-card>
            <v-card-title  class="headline font-weight-thin">
              Visites d'entreprises
            </v-card-title>
            <v-divider/>
            <v-card-text>
              <v-list>
                <v-list-tile avatar>
                  Planifiées: x%
                </v-list-tile>
                <v-list-tile avatar>
                  Effectuées: y%
                </v-list-tile>
                <v-list-tile avatar>
                  Retard: y%
                </v-list-tile>
              </v-list>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
  </v-container>

  </div>
</template>

<script>
export default {
  computed: {
    message () {
      return this.$store.getters.reverseLog
    },
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    }
  }
}
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
div.titre {
  font-family: 'Signika', sans-serif;
  color: #ffffff;
  font-weight: 500;
  font-size: 18px
}
</style>
