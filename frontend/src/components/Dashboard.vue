<template>
  <div>
    <v-toolbar class="toolbar elevation-12" color="#c9aec5" height="35px" app>
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
      <v-icon color="#c9aec5" v-if="!rightDrawer" @click="rightDrawer=!rightDrawer">fa-dashboard</v-icon>
    </v-toolbar>
    <v-navigation-drawer
      :class="rightDrawer?'elevation-6':''"
      right app
      v-model="rightDrawer"
    >
      <v-list  two-line class="pt-0">
        <v-toolbar>
          <v-icon @click="rightDrawer=!rightDrawer" color="c9aec5">fa-dashbord</v-icon>
          <v-spacer></v-spacer>
            Suivi d'activité
          <v-divider></v-divider>
        </v-toolbar>
      </v-list>
    </v-navigation-drawer>
    <v-card v-for="t in tasks" :key="t._id" class="ma-3 task elevation-12">
      <v-card-text class="tasktitle">
        <v-tooltip>
          <span>SIRET: {{ t._id }}</span>
          <v-card-title slot="activator">{{ t.etablissement[0].value.sirene.raisonsociale }}</v-card-title>
        </v-tooltip>
      </v-card-text>
      <v-card-text class="tasktext" dark v-for="s in t.tasks" :key="s.id" style="background-color: transparent;">
        {{ s }}
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import Etablissement from '@/components/Etablissement'

export default {
  components: { Etablissement },
  data () {
    return {
      tasks: []
    }
  },
  watch: {
    search (val) {
      val && val !== this.select && this.querySelections(val)
    }
  },
  methods: {
    getTasks () {
      this.$axios.get('/api/dashboard/tasks').then(response => {
        this.tasks = response.data
      })
    }
  },
  mounted () {
    this.getTasks()
  },
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
    },
    rightDrawer: {
      get () {
        return this.$store.state.rightDrawer
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
  }
}
</script>

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
