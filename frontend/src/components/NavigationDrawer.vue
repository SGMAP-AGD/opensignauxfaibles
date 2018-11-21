<template>
  <v-navigation-drawer v-model="drawer" app>
    <v-toolbar flat class="transparent">
      <v-list class="pa-0">
        <v-list-tile avatar>
          <v-list-tile-avatar>
            <img src="/static/logo_signaux_faibles_cercle.svg">
          </v-list-tile-avatar>

          <v-list-tile-content>
            <v-list-tile-title>Signaux·Faibles</v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-avatar>
            <v-icon @click="logout()">logout</v-icon>
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
      </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  methods: {
    logout () {
      this.$store.commit('logout')
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
    }
  }
}
</script>