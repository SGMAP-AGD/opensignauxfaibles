import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'

Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

const state = {
  token: ''
}

const mutations = {
  setToken (state, token) {
    state.token = token
  }
}

export default new Vuex.Store({
  state,
  mutations,
  plugins: [vuexLocal.plugin]
})
