<template>
  <v-app id="app">
    <v-content>
        <NavigationDrawer v-if="login && drawer"/>
        <Login v-if="!login" :state="login"/>
        <router-view v-if="login" />
    </v-content>
  </v-app>
</template>

<script>
import Login from '@/components/Login'
import Uploads from '@/components/Uploads'
import NavigationDrawer from '@/components/NavigationDrawer'

export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    }
  },
  components: { Login, Uploads, NavigationDrawer },
  computed: {
    uploads () {
      return this.$store.getters.getUploads
    },
    login () {
      return this.$store.state.token != null
    },
    batches () {
      return JSON.stringify(this.$store.state.batches, null, 2)
    },
    dbstatus () {
      return this.$store.state.dbstatus
    },
    messages () {
      return this.$store.getters.messages
    },
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    }
  },
  data () {
    return {
      fixed: false,
      dialog: false,
      activityFilter: ['info', 'warning', 'critical'],
      menu: {
        title: 'Accueil',
        color: 'light-green darken-4'
      },
      items: [
        { title: 'Accueil',
          action: '/',
          icon: 'fa-home',
          color: 'green darken-3' },
        { title: 'Détection',
          action: '/Browse',
          icon: 'fa-search',
          color: 'indigo darken-4' },
        { title: 'Données',
          action: '/data',
          icon: 'fa-database',
          color: 'red darken-4' },
        { title: 'Administration',
          action: '/admin',
          icon: 'fa-users',
          color: 'gray darken-4' }
      ],
      avatar: {
        debug: {icon: 'fa-cogs', color: 'blue'},
        info: {icon: 'fa-info', color: 'green'},
        warning: {icon: 'fa-exclamation-triangle', color: 'yellow'},
        critical: {icon: 'fa-sad-cry', color: 'red'}
      }
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

<style>
 @import url('https://fonts.googleapis.com/css?family=Quicksand');
 @import url('https://fonts.googleapis.com/css?family=Signika');
  body {
    font-family: 'Quicksand', sans-serif;
  }
  .toolbar {
    background:         radial-gradient( circle at center, red, blue);
    color: "black";
    font-family: 'Quicksand', sans-serif;
    background-repeat: repeat;
    background-image: url("/static/bgtoolbar.png");
  }
  #app {
    background-color: #e7e7e7;
    background: radial-gradient(circle at center, rgb(255, 255, 255), rgb(228, 228, 228) 75%, rgb(204, 204, 204) 100%);
  }
  .rightDrawer {
    position: fixed;
    right: 0px;
  }
  .span {
    max-height: 10px
}
</style>