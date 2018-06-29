<template>
  <v-app>
    <v-toolbar
      app
      :clipped-left="clipped"
    >
    <v-menu open-on-hover top offset-y>
      <v-btn slot="activator" width="100px" :color="menu.color" dark>{{ menu.title }}</v-btn>
      <v-list>
        <v-list-tile v-for="(item, index) in items" 
        :key="index" 
        @click="setMenu(index)"
        :to="item.action"
        :color="item.color">
          <v-list-tile-title>{{ item.title }}</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
    <v-toolbar-title>
      <span id="gray">Open</span>
      <span id="blue">Signaux</span>
      <span id="red">Faibles</span>
    </v-toolbar-title>
    <v-spacer></v-spacer>
    <v-btn icon>
      <v-icon>fa-user</v-icon>
    </v-btn>
    </v-toolbar>
    <v-content>
      <router-view/>
    </v-content>

    <v-footer>
      <v-btn 
      flat 
      icon 
      color="blue"
      href="https://github.com/entrepreneur-interet-general/opensignauxfaibles">
        <v-icon >fa-github</v-icon>
      </v-btn>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    }
  },
  data () {
    return {
      clipped: false,
      fixed: false,
      drawer: false,
      menu: {
        title: 'Menu',
        color: 'light-green darken-4'
      },
      items: [
        { title: 'Détection',
          action: '/prediction',
          color: 'indigo darken-4' },
        { title: 'Données',
          action: '/data',
          color: 'grey darken-2' },
        { title: 'Visites',
          action: '/workflow',
          color: 'red darken-4' },
        { title: 'Admin',
          action: '/admin',
          color: 'black' }
      ]
    }
  },
  name: 'App'
}
</script>

<style scoped>
 @import url('https://fonts.googleapis.com/css?family=Quicksand');
body {
  font-family: 'Quicksand', sans-serif;
}
#gray {
  color: #606060;
  font-family: 'Quicksand', sans-serif;
}
#blue {
  color: #20449a;
  font-family: 'Quicksand', sans-serif;
}
#red {
  color: #e9222e;
  font-family: 'Quicksand', sans-serif;
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
v-list {
  font-family: 'Quicksand', sans-serif;
}
</style>

