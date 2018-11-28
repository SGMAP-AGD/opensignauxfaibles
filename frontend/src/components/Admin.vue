<template>
<div>
  <v-toolbar class="toolbar" height="35px" color="#ffffff" dense app>
    <v-icon
      class="fa-rotate-180"
      @click="drawer=!drawer"
      v-if="!drawer"
      color="#444"
      key="toolbar"
    >mdi-backburger</v-icon>
    <div style="width: 100%; text-align: center;"  class="titre">
      Administration
    </div>
    <v-spacer></v-spacer>
    <v-icon color="#444"  @click="rightDrawer=!rightDrawer">fa-cog</v-icon>
  </v-toolbar>
  <v-container grid-list-md text-xs-center>
    <v-layout>
      <v-flex key="database" xs3>
        <v-card>
          <v-card-title databaseCopy>
            Copie de la base
          </v-card-title>
          <v-card-actions>
            <v-text-field
              name="to"
              label="Destination"
              id="db"
              v-model="to"
            ></v-text-field>
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="cloneDatabase()">Cloner</v-btn>
          </v-card-actions>
        </v-card>
        <v-card>
          <v-card-title databaseCopy>
            Purge de la base
          </v-card-title>
          <v-card-actions>
            <v-btn @click="purgeDatabase()">Purger</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
      <v-flex xs3>
        <v-btn @click="go()">go</v-btn>
      </v-flex>
      <v-flex xs3>
        {{ d }}
      </v-flex>
    </v-layout>
  </v-container>
</div>
</template>

<script>
import axios from 'axios'

export default {
  methods: {
    cloneDatabase () {
      axios.get('http://localhost:3000/api/admin/clone/' + this.to)
        .then(response => alert(JSON.stringify(response.data, null, 2)))
    },
    go () {
      this.$axios.get('/api/processBatch').then(r => { this.d = r.data })
    },
    purgeDatabase () {
      this.$axios.get('/api/purge')
        .then(response => alert(JSON.stringify(response.data, null, 2)))
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
  },
  data () {
    return {
      to: '',
      d: null
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
.drawer {
  z-index: 1;
}
div.titre {
  color: #444;
  font-family: 'Signika', sans-serif;
  font-weight: 500;
  color: primary;
    font-size: 18px;
}
</style>
