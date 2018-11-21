<template>
  <v-container fluid fill-height>
    <v-layout align-center wrap justify-center>
      <v-flex class="login" xs12 sm6 md3>
        <span><v-img class="center" max-width="75px" src="/static/logo_signaux_faibles_cercle.svg"></v-img></span>
        <h1><span class="fblue">Signaux</span>·<span class="fred">Faibles</span></h1>
        <hline></hline>
      </v-flex>
      <v-flex class="login" xs12 sm6 md3>
          <v-form @submit="login">
          <v-text-field solo prepend-icon="person" label="Nom d'utilisateur" v-model="username"></v-text-field >
          <v-text-field solo prepend-icon="lock" type="password" label="Mot de passe" v-model="password"></v-text-field>
          <v-btn color="primary" type="submit">Login</v-btn>
          </v-form>
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

<style scoped>
  .login {
        text-align: center;     /* will center text in <p>, which is not a flex item */
        align-items: center
  }
  h1 {
    font-family: 'Quicksand', sans-serif;
  }
  .center {
    display: block;
    margin: 0 auto;
  }
  #toolbar {
    font-size: 25px;
    color: white;
  }
  span.fblue {
      color: #20459a
  }
  span.fred {
      color: #e9222e
  }
  .box {
    background-color: #fff;
  }
</style>