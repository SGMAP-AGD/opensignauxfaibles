import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

var axiosClient = axios.create(
  {
    headers: {
      'Content-Type': 'application/json'
    },
    baseURL: 'http://localhost:3000'
  }
)

axiosClient.interceptors.request.use(
  config => {
    if (store.state.token != null) config.headers['Authorization'] = 'Bearer ' + store.state.token
    return config
  }
)

const store = new Vuex.Store({
  plugins: [createPersistedState({storage: window.sessionStorage})],
  state: {
    credentials: {
      username: null,
      password: null
    },
    token: null,
    types: null,
    features: null,
    files: null,
    batches: [],
    dbstatus: null,
    currentBatchKey: 0,
    lastMove: 0
  },
  mutations: {
    login (state) {
      axiosClient.post('/login', state.credentials).then(response => {
        state.token = response.data.token
        store.commit('updateRefs')
        store.commit('updateBatches')
        store.commit('updateDbStatus')
      })
    },
    refreshToken (state) {
      axiosClient.get('/api/refreshToken').then(response => {
        state.token = response.data.token
      })
    },
    logout (state) {
      state.credentials.username = null
      state.credentials.password = null
      state.token = null
      state.types = null
      state.features = null
      state.files = null
      state.batches = []
      state.lastMove = 0
    },
    setUser (state, username) {
      state.credentials.username = username
    },
    setPassword (state, password) {
      state.credentials.password = password
    },
    updateBatches (state) {
      axiosClient.get('/api/admin/batch').then(response => {
        state.batches = response.data
      })
    },
    updateDbStatus (state) {
      axiosClient.get('/api/admin/status').then(response => {
        state.dbstatus = response.data
      })
    },
    setLastMove (state, lastMove) {
      state.lastMove = lastMove
    },
    updateRefs (state) {
      axiosClient.get('/api/admin/types').then(response => { state.types = response.data.sort() })
      axiosClient.get('/api/admin/features').then(response => { state.features = response.data })
      axiosClient.get('/api/admin/files').then(response => { state.files = response.data })
    },
    setCurrentBatchKey (state, key) {
      state.currentBatchKey = key
    }
  },
  getters: {
    axiosConfig (state) {
      return {headers: {Authorization: 'Bearer ' + state.token}}
    }
  }
})

setInterval(
  function () {
    if (store.state.token != null) {
      axiosClient.get('/api/lastMove').then(response => {
        if (response.data > store.state.lastMove) {
          store.commit('setLastMove', response.data)
          store.commit('updateRefs')
          store.commit('updateBatches')
          store.commit('updateDbStatus')
        }
      }).catch(error => {
        if (error.response.status === 401) {
          store.commit('logout')
        }
        console.log(error.response)
      })
    }
  },
  500)

if (store.state.token != null) {
  store.commit('refreshToken')
}

setInterval(
    function () {
      if (store.state.token != null) {
        store.commit('refreshToken')
      }
    },
    30000)

store.commit('updateRefs')
store.commit('updateDbStatus')
store.commit('updateBatches')

export default store
