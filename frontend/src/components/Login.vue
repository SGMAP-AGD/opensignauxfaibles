<template>
<v-container align-center>
  <v-layout center>
    <v-flex xs4 offset-xs4>
      <v-card>
        <v-toolbar card color="light-blue lighten-2">
          Authentification
        </v-toolbar>
        <v-card-text>
          <v-text-field label="Nom d'utilisateur" v-model="user">

          </v-text-field >
        <br/>
          <v-text-field type="password" label="Mot de passe" v-model="password">

          </v-text-field>
        </v-card-text>
        <v-card-action>
          <v-btn color="blue" flat @click="login()">LOGIN</v-btn>
        </v-card-action>
      </v-card>
    </v-flex>
  </v-layout>
</v-container>
</template>

<script>
export default {
  methods: {
    login () {
      if (this.user === 'test' && this.password === 'test') {
        this.state.state = true
      }
    }
  },
  props: {
    state: {
      type: Object,
      default: () => ({})
    }
  },
  data () {
    return {
      user: '',
      password: ''
    }
  }
}
</script>