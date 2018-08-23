import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

const store = new Vuex.Store({
  plugins: [createPersistedState({storage: window.sessionStorage})],
  state: {
    api: 'http://localhost:3000/',
    credentials: {
      username: null,
      password: null
    },
    token: null
  },
  mutations: {
    login (state) {
      axios.post(state.api + 'login', state.credentials).then(response => { state.token = response.data.token })
    },
    logout (state) {
      state.credentials.username = null
      state.credentials.password = null
      state.token = null
    },
    setUser (state, username) {
      state.credentials.username = username
    },
    setPassword (state, password) {
      state.credentials.password = password
    }
  },
  getters: {
    axiosConfig (state) {
      return {headers: {Authorization: 'Bearer ' + state.token}}
    }
  }
})

export default store
