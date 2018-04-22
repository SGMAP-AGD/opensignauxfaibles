<template>
  <div>
    <md-dialog id="login" :md-active.sync="showDialog">
      <md-card>
        <md-card-header>
          <md-dialog-title>Connexion</md-dialog-title>
        </md-card-header>
        <md-card-content>
          <form @submit.prevent="login">
            <md-field>
              <label for="Nom d'utilisateur">Utilisateur</label>
              <md-input name="username" id="username" v-model="username"/>
            </md-field>
            <md-field>
              <label for="Mot de passe">Mot de passe</label>
              <md-input type="password" name="password" id="password" v-model="password"/>
            </md-field>
            <md-button type="submit" class="md-primary" >Connexion</md-button>
            <md-button class="md-primary" @click="showDialog = false">Fermer</md-button>
          </form>
        </md-card-content>
      </md-card>
    </md-dialog>
    <md-button v-if="$store.state.token==''" class="md-primary md-raised" @click="showDialog = true">Se connecter</md-button>
    <div v-if="$store.state.token!=''">
      <md-button class="md-primary md-raised" @click="logout">Se déconnecter</md-button>
    </div>
  </div>

</template>
<script>
import axios from 'axios'

export default {
  name: 'Login',
  data: () => ({
    showDialog: false,
    username: '',
    password: '',
    token: ''
  }),
  methods: {
    login () {
      axios.get(
        `http://localhost:3000/api/auth`,
        {
          auth: {
            username: this.username,
            password: this.password
          }
        }
      ).then((response) => {
        if (response.status === 200) {
          this.showDialog = false
          this.$store.commit('setToken', response.data.token)
        }
      })
    },
    logout () {
      this.username = ''
      this.password = ''
      this.$store.commit('setToken', '')
    }
  }
}
</script>

<style>
  .md-dialog {
    width: 250px;
  }
  .md-field {
    width: 85%px
  }
</style>
