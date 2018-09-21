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
    files: [],
    batches: [],
    dbstatus: null,
    currentBatchKey: 0,
    currentType: null,
    epoch: 0,
    socket: {
      isConnected: false,
      message: [],
      reconnectError: false
    }
  },
  mutations: {
    SOCKET_ONOPEN (state, event) {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state, event) {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event) {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message) {
      state.socket.message.push(message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true
    },
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
      state.epoch = 0
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
    setEpoch (state, epoch) {
      state.epoch = epoch
    },
    updateRefs (state) {
      axiosClient.get('/api/admin/types').then(response => { state.types = response.data.sort((a, b) => a.text.localeCompare(b.text)) })
      axiosClient.get('/api/admin/features').then(response => { state.features = response.data })
      axiosClient.get('/api/admin/files').then(response => { state.files = response.data.sort((a, b) => a.name.localeCompare(b.name)) })
    },
    setCurrentBatchKey (state, key) {
      state.currentBatchKey = key
    },
    setCurrentType (state, type) {
      state.currentType = type
    }
  },
  actions: {
    saveBatch (state, batch) {
      console.log(batch)
      axiosClient.post('/api/admin/batch', batch).then(r => { state.currentBatch = batch })
    },
    checkEpoch () {
      if (store.state.token != null) {
        axiosClient.get('/api/admin/epoch').then(response => {
          if (response.data !== store.state.epoch) {
            store.commit('setEpoch', response.data)
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
    }
  },
  getters: {
    axiosConfig (state) {
      return {headers: {Authorization: 'Bearer ' + state.token}}
    }
  }
})

// setInterval(
//   function () {
//     store.dispatch('checkEpoch')
//   },
//   500)

if (store.state.token != null) {
  store.commit('refreshToken')
}

setInterval(
    function () {
      if (store.state.token != null) {
        store.commit('refreshToken')
      }
    },
    180000)

store.commit('updateRefs')
store.commit('updateDbStatus')
store.commit('updateBatches')

export default store
