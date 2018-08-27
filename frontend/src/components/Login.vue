<template>
<v-container align-center>
  <v-layout center>
    <v-flex xs4 offset-xs4>
      <v-card>
        <v-toolbar class="toolbar" id="toolbar" card color="indigo darken-1">
          <v-spacer> </v-spacer>
            Identification
          <v-spacer> </v-spacer>
        </v-toolbar>
        <v-card-text>
          <v-text-field label="Nom d'utilisateur" v-model="username" outline></v-text-field >
          <v-text-field type="password" label="Mot de passe" v-model="password" outline></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-btn color="blue" flat @click="login">LOGIN</v-btn>
        </v-card-actions>
      </v-card>
    </v-flex>
  </v-layout>
</v-container>
</template>

<script>

export default {
  computed: {
    username: {
      get (username) {
        return this.$store.state.credentials.username
      },
      set (username) {
        this.$store.commit('setUser', username)
      }
    },
    password: {
      get (password) {
        return this.$store.state.credentials.password
      },
      set (password) {
        this.$store.commit('setPassword', password)
      }
    },
    token () {
      return this.$store.state.token
    }
  },
  methods: {
    login () {
      this.$store.commit('login')
    }
  }
}
</script>

<style>
  #toolbar {
    font-size: 25px;
    color: white;
  }
</style>