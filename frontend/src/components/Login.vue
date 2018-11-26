<template>
  <v-container fluid fill-height>
    <v-layout align-center wrap justify-center>
      <v-flex class="login" xs12 sm6 md3>
        <span><v-img class="center" max-width="100px" src="/static/logo_signaux_faibles_cercle.svg"></v-img></span>
        <h1><span class="fblue">Signaux</span>·<span class="fred">Faibles</span></h1>
      </v-flex>
      <v-flex class="login" xs12 sm6 md3>
          <v-form @submit="login">
          <v-text-field solo prepend-icon="person" label="Adresse électronique" v-model="username"></v-text-field >
          <v-text-field solo prepend-icon="lock" type="password" label="Mot de passe" v-model="password"></v-text-field>
          <v-btn round color="primary" type="submit">Connexion</v-btn><br/>
          <v-dialog 
            height="400px"
            width="500"
            v-model="passwordHelp"
          >
          <a slot="activator" style="font-size: 10px" href='#'>Mot de passe oublié ?</a>
          <v-card >
            <v-toolbar color="#eef" class="toolbar elevation-1" dense>
              <v-toolbar-avatar><v-icon>mdi-lock-question</v-icon></v-toolbar-avatar>
              <v-spacer></v-spacer>
              <span style="font-size: 15px; font-weight: 600">
              Récupération du compte
              </span>
            </v-toolbar>

            <div v-if="state==1">
              <v-card-text>
                <b>Etape 1</b> · Etape 2 · Fin<br/>
                Saisissez votre adresse électronique, un code de récupération vous sera envoyé.
              </v-card-text>
              <v-card-text>
                <v-text-field  prepend-icon="mdi-at" label="Adresse E-Mail" v-model="recovery"></v-text-field >
              </v-card-text>
                <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="primary"
                  flat
                  @click="getRecovery()"
                >
                  suivant
                </v-btn>
              </v-card-actions>
            </div>
            <div v-if="state==2">
              <v-card-text>
                Etape 1 · <b>Etape 2</b> · Fin<br/>
                Saisissez le code de récupération reçu dans votre boite aux lettres électronique
              </v-card-text>
              <v-card-text>
                <v-text-field length=6 solo prepend-icon="mdi-pound" label="Code de récupération" v-model="recoveryCode"></v-text-field >
              </v-card-text>
                <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="primary"
                  flat
                  @click="checkRecovery()"
                >
                  suivant
                </v-btn>
              </v-card-actions>
            </div>
          </v-card>
          </v-dialog>
          </v-form>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>

export default {
  data () {
    return {
      passwordHelp: false,
      recoveryEmail: null,
      recoveryCode: null,
      state: 1
    }
  },
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
    },
    getRecovery () {
      let self = this
      this.$axios.get('/recovery/get/' + this.pass).then(r => {
        self.state = 2
      })
    },
    checkRecovery () {
      let self = this
      this.$axios.get('/recovery/check/' + this.recovery + '/' + this.recoveryCode).then(r => {
        self.state = 2
      })
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
    font-family: 'Quicksand', sans-serif;
    color: #20459a
  }
  span.fred {
    font-family: 'Quicksand', sans-serif;
    color: #e9222e
  }
  .box {
    background-color: #fff;
  }
</style>