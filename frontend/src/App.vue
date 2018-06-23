<template>
  <v-app>
    <v-navigation-drawer
      persistent
      :mini-variant="miniVariant"
      clipped
      v-model="drawer"
      permanent
      fixed
      app
    >
      <v-list>
        <v-list-tile
          value="true"
          v-bind:to="item.to"
          v-for="(item, i) in items"
          :key="i"
        >
          <v-list-tile-action>
            <v-icon v-html="item.icon"></v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title v-text="item.title"></v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar
      app
      :clipped-left="clipped"
    >
      <v-btn icon @click.stop="miniVariant = !miniVariant">
        <v-icon v-html="miniVariant ? 'chevron_right' : 'chevron_left'"></v-icon>
      </v-btn>


      <v-toolbar-title>Open <span id="blue">Signaux</span> <span id="red">Faibles</span></v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon>fa-user</v-icon>
      </v-btn>
    </v-toolbar>
    <v-content>
      <router-view/>
    </v-content>

    <v-footer :fixed="fixed" app>
      <span>&copy; 2017</span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  data () {
    return {
      clipped: false,
      fixed: true,
      items: [{
        icon: 'fa-heartbeat',
        title: 'Prédictions',
        to: '/prediction'
      }, {
        icon: 'fa-building',
        title: 'Entreprises',
        to: '/'
      }, {
        icon: 'fa-cogs',
        title: 'Gestion des données',
        to: '/data'
      }, {
        icon: 'fa-bug',
        title: 'Signaler un problème',
        to: '/bug'
      }],
      miniVariant: false,
      right: true,
      rightDrawer: false
    }
  },
  name: 'App'
}
</script>

<style scoped>
       #blue {
         color: #20449a;
       }
      #red {
         color: #e9222e;
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
</style>

