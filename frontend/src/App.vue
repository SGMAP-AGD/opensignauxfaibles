<template>
  <v-app id="app">
    <v-toolbar
    class="toolbar"
    color="blue-grey lighten-4"
    app
    >
    <v-menu open-on-hover top offset-y transition="slide-x-transition">
      <v-btn  slot="activator" width="100px" :color="menu.color" dark>{{ menu.title }}</v-btn>
      <v-list>
        <v-list-tile
        v-for="(item, index) in items"
        :key="index"
        @click="setMenu(index)"
        :to="item.action"
        :color="item.color"
          style="vertical-align: middle">
          <v-list-tile-title
          style="vertical-align: middle">
            {{ item.title }}
          </v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
    <v-toolbar-title>
      <span id="blue">Signaux</span>
      <span id="red">Faibles</span><br/>
      <span class="caption gray">Plateforme de détection des entreprises fragiles</span>
    </v-toolbar-title>
    <v-spacer></v-spacer>

      <v-btn flat icon disabled>
      <v-tooltip bottom open-delay="1000" transition="fade-transition" v-if="login && dbstatus != null">
        {{ dbstatus }}
          <v-progress-circular slot="activator"
            indeterminate
            color="indigo darken-4"
          ></v-progress-circular>
      </v-tooltip>

      <!-- <v-tooltip bottom open-delay="1000" transition="fade-transition" v-if="login && dbstatus === null">
        Aucune opération en cours
          <v-icon slot="activator" color="green">fa-check
          </v-icon>
      </v-tooltip> -->
      </v-btn>

      <v-dialog 
      lazy 
      v-model="dialog" 
      scrollable v-if="login">
        <v-btn flat lazy icon slot="activator" >
          <v-icon>fa-newspaper</v-icon>
        </v-btn>
        <v-card
        style="min-height: 90vh;">
          <v-toolbar
          class="headline"
          color="indigo darken-4"
          dark
          dense
          card>
            <v-toolbar-title>
              Activité Serveur
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <v-toolbar-items>
              <v-switch
              height="100%"
              v-for="(key,index) in avatar"
              :key="key.color"
              v-model="activityFilter" 
              :color="key.color"
              :prependIcon="key.icon" 
              :value="index"/>
              <v-spacer></v-spacer>
              <v-icon color="red" @click="dialog = false">fa-window-close</v-icon>
              </v-toolbar-items>
            </v-toolbar>
          <v-card-text >
            <v-list dense>
              <template 
              v-for="(m, index) in messages.filter(m => activityFilter.indexOf(m.priority) > -1)">
                <v-divider :key="'divide' + index"/>
                <v-list-tile
                  :key="index"
                  avatar
                >
                  <v-list-tile-avatar>
                    <v-icon :color="(avatar[m.priority]||{}).color">{{ (avatar[m.priority]||{}).icon }}</v-icon>
                  </v-list-tile-avatar> 

                  <v-list-tile-content>
                    <v-list-tile-title v-html="m.event"/>
                    <v-list-tile-sub-title class="caption" v-html="m.date.toLocaleString()"/>
                  </v-list-tile-content>
                </v-list-tile>
              </template>
            </v-list>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-tooltip bottom open-delay="1000" transition="fade-transition" v-if="login" >
      Déconnexion
          <v-btn icon @click="logout()" slot="activator">
            <v-icon>fa-sign-out-alt</v-icon>
          </v-btn>
      </v-tooltip>

    </v-toolbar>
    <v-content>
        <Login v-if="!login" :state="login"/>
        <router-view v-if="login" />
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
      this.$store.commit('logout')
    }
  },
  components: { Login },
  computed: {
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
    }
  },
  data () {
    return {
      fixed: false,
      drawer: false,
      dialog: false,
      activityFilter: ['info', 'warning', 'critical'],
      menu: {
        title: 'Accueil',
        color: 'light-green darken-4'
      },
      items: [
        { title: 'Accueil',
          action: '/',
          color: 'green darken-3' },
        { title: 'Détection',
          action: '/prediction',
          color: 'indigo darken-4' },
        { title: 'Données',
          action: '/data',
          color: 'grey darken-2' }
        // { title: 'Admin',
        //   action: '/admin',
        //   color: 'black'},
        // { title: 'Test',
        //   action: '/admintoto',
        //   color: 'yellow'}
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

<style scoped>
 @import url('https://fonts.googleapis.com/css?family=Quicksand');
body {
  font-family: 'Quicksand', sans-serif;
}
.gray {
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

<style>
.toolbar {
  background-image: url("/static/bgtoolbar.png");
  color: "black";
}
#app {
  background-image: url("/static/bgapp.png");
  background-color: "blue-grey lighten-5";
}
</style>