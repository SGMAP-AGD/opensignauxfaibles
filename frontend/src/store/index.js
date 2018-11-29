import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import createPersistedState from 'vuex-persistedstate'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)

const vm = new Vue()

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
    if (sessionStore.state.token != null) config.headers['Authorization'] = 'Bearer ' + sessionStore.state.token
    return config
  }
)

const localStore = new Vuex.Store({
  plugins: [createPersistedState({ storage: window.localStorage })],
  state: {
    browserToken: null
  },
  mutations: {
    setBrowserToken (state, browserToken) {
      state.browserToken = browserToken
    }
  },
  getters: {
    browserToken (state) { return state.browserToken }
  }
})

const sessionStore = new Vuex.Store({
  plugins: [createPersistedState({ storage: window.sessionStorage })],
  state: {
    credentials: {
      email: null,
      password: null
    },
    appDrawer: false,
    rightDrawer: false,
    token: null,
    types: null,
    features: null,
    files: [],
    batches: [],
    dbstatus: null,
    currentBatchKey: null,
    currentType: null,
    epoch: 0,
    socket: {
      isConnected: false,
      message: [],
      reconnectError: false
    },
    uploads: [],
    activeTab: null,
    tabs: [{
      'type': 'Prediction',
      'param': '09/2018',
      'batch': '1802'
    }
    ]
  },
  mutations: {
    updateActiveTab (state, activeTab) {
      state.activeTab = activeTab
    },
    drawer (state, val) {
      state.appDrawer = val
    },
    rightDrawer (state, val) {
      state.rightDrawer = val
    },
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
    SOCKET_ONMESSAGE (state, message) {
      if ('journalEvent' in message) {
        let m = message.journalEvent
        state.socket.message.unshift(m)
        if (state.socket.message.length > 250) {
          state.socket.message.pop()
        }
      }
      if ('batches' in message) {
        state.batches = message.batches.reverse()
      }
      if ('files' in message) {
        state.files = message.files
      }
    },
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      wsConnect(state)
    },

    refreshToken (state) {
      axiosClient.get('/api/refreshToken').then(response => {
        state.token = response.data.token
      })
    },
    logout (state) {
      state.credentials.email = null
      state.credentials.password = null
      state.token = null
      state.types = null
      state.features = null
      state.files = null
      state.batches = []
      state.epoch = 0
      state.socket.message = []
      state.uploads = []
    },
    setEmail (state, email) {
      state.credentials.email = email
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
    updateLogs (state) {
      axiosClient.get('/api/admin/getLogs').then(response => { state.socket.message = (response.data || []) })
    },
    setCurrentBatchKey (state, key) {
      state.currentBatchKey = key
    },
    setCurrentType (state, type) {
      state.currentType = type
    },
    updateUploads (state, status) {
      let index = state.uploads.findIndex(s => s.name === status.name)
      state.uploads[index].amount = status.amount
    },
    addUpload (state, status) {
      if (state.uploads.findIndex(s => s.name === status.name) === -1) state.uploads.push(status)
    },
    resetUploads (state) {
      state.uploads = state.uploads.filter(u => u.amount < 100)
    },
    updateTabs (state, tabs) {
      state.tabs = tabs
    },
    setToken (state, token) {
      state.token = token
    }
  },
  actions: {
    updateBatches (context) {
      axiosClient.get('/api/admin/batch').then(response => {
        context.commit('updateBatches', response.data)
      })
    },
    setCurrentType (context, type) {
      context.commit('setCurrentType', type)
    },
    login (context) {
      let credentials = {
        email: context.state.credentials.email,
        password: context.state.credentials.password,
        browserToken: localStore.state.browserToken
      }

      axiosClient.post('/login', credentials).then(response => {
        context.commit('setToken', response.data.token)
        wsConnect(context)
        context.commit('updateRefs')
        context.commit('updateBatches')
        context.commit('updateDbStatus')
        context.commit('updateLogs')
      })
    },
    getLogin (context) {
      let credentials = {
        email: context.state.credentials.email,
        password: context.state.credentials.password
      }

      axiosClient.post('/login/get', credentials)
    },
    checkLogin (context, checkCode) {
      let credentials = {
        email: context.state.credentials.email,
        password: context.state.credentials.password,
        checkCode: checkCode
      }

      axiosClient.post('/login/check', credentials).then(response => {
        console.log(response.data)
        localStore.commit('setBrowserToken', response.data.browserToken)
        context.commit('login')
      })
    },
    setDrawer (context, val) {
      context.commit('drawer', val)
    },
    setRightDrawer (context, val) {
      context.commit('rightDrawer', val)
    },
    updateActiveTab (context, activeTab) {
      context.commit('updateActiveTab', activeTab)
    },
    updateTabs (state, tabs) {
      state.commit('updateTabs', tabs)
    },
    saveBatch (state, batch) {
      axiosClient.post('/api/admin/batch', batch).then(r => { state.currentBatch = batch })
    },
    resetUploads (context) {
      context.commit('resetUploads')
    },
    addFile (context, file) {
      axiosClient.post('/api/admin/batch/addFile', file)
    },
    upload (context, file) {
      let formData = new FormData()
      let filename = '/' + file.batch + '/' + file.type + '/' + file.file.name

      formData.append('file', file.file)
      formData.append('batch', file.batch)
      formData.append('type', file.type)

      var status = {
        'amount': 0,
        'name': filename
      }

      context.commit('addUpload', status)
      axiosClient.post(
        '/api/admin/files',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          onUploadProgress: function (progressEvent) {
            var newStatus = {
              'amount': parseInt(Math.round((progressEvent.loaded * 100) / progressEvent.total)),
              'name': filename
            }
            context.commit('updateUploads', newStatus)
            context.commit('updateBatches')
          }
        }
      ).then(response => {
        let postData = {
          filename: filename,
          type: file.type,
          batch: file.batch
        }
        context.dispatch('addFile', postData)
      }).catch(function (response) {
      })
    },
    checkEpoch (context) {
      if (context.state.token != null) {
        axiosClient.get('/api/admin/epoch').then(response => {
          if (response.data !== sessionStore.state.epoch) {
            context.commit('setEpoch', response.data)
            context.commit('updateRefs')
            context.commit('updateBatches')
            context.commit('updateDbStatus')
          }
        }).catch(error => {
          if (error.response.status === 401) {
            context.commit('logout')
          }
        })
      }
    }
  },
  getters: {
    batches (state) {
      return state.batches.reduce((accu, batch) => {
        accu[batch.id.key] = batch
        return accu
      }, {})
    },
    axiosConfig (state) {
      return { headers: { Authorization: 'Bearer ' + state.token } }
    },
    messages (state) {
      return state.socket.message.map(m => {
        m.date = new Date(m.date)
        return m
      })
    },
    getUploads (state) {
      return state.uploads
    },
    getTabs (state) {
      return state.tabs
    },
    activeTab (state) {
      return state.activeTab
    }
  }
})

function wsConnect (state) {
  let index = Vue._installedPlugins.indexOf(VueNativeSock)
  if (index > -1) {
    Vue._installedPlugins.splice(index, 1)
  }
  Vue.use(VueNativeSock, 'ws://localhost:3000/ws/' + state.token, {
    store: sessionStore,
    format: 'json',
    connectManually: true,
    reconnection: true,
    reconnectionAttempts: -1,
    reconnectionDelay: 3000
  })
  vm.$connect()
}

if (sessionStore.state.token != null) {
  wsConnect(sessionStore.state)
  sessionStore.commit('refreshToken')
  sessionStore.commit('updateRefs')
  sessionStore.commit('updateDbStatus')
  sessionStore.commit('updateBatches')
  sessionStore.commit('updateLogs')
}

setInterval(
  function () {
    if (sessionStore.state.token != null) {
      sessionStore.commit('refreshToken')
    }
  },
  180000)

var store = {
  sessionStore: sessionStore,
  localStore: localStore
}

export default store
