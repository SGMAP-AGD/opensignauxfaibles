<template>
  <v-app>
    <v-toolbar
      app
      :clipped-left="clipped"
    >
    <v-menu open-on-hover top offset-y transition="slide-x-transition">
      <v-btn  slot="activator" width="100px" :color="menu.color" dark>{{ menu.title }}</v-btn>
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
    
    <v-btn icon @click="logout()">
      <img height="28px" src="/static/favicon.png" />
    </v-btn>
    </v-toolbar>
    <v-content>
      <Login v-if="!login.state" :state="login"/>
      <router-view v-if="login.state" />
    </v-content>

    <v-footer>
      <v-btn 
        flat 
        icon 
        color="blue"
        href="https://github.com/entrepreneur-interet-general/opensignauxfaibles">
        <v-icon>fab fa-github</v-icon>
      </v-btn>
    </v-footer>
  </v-app>
</template>

<script>
import Login from '@/components/Login'

export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    },
    logout () {
      this.login.state = false
    }
  },
  components: { Login },
  data () {
    return {
      login: {state: false},
      clipped: false,
      fixed: false,
      drawer: false,
      username: '',
      password: '',
      menu: {
        title: 'Accueil',
        color: 'light-green darken-4'
      },
      items: [
        { title: 'Accueil',
          action: '/',
          color: 'green' },
        { title: 'Détection',
          action: '/prediction',
          color: 'indigo darken-4' },
        { title: 'Données',
          action: '/data',
          color: 'grey darken-2' }
        //   ,
        // { title: 'Visites',
        //   action: '/workflow',
        //   color: 'red darken-4' },
        // { title: 'Admin',
        //   action: '/admin',
        //   color: 'black' }
      ]
    }
  },
  mounted () {
    this.items.forEach((item, index) => {
      if (item.action === location.hash.substring(1)) {
        this.setMenu(index)
      }
    })
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

