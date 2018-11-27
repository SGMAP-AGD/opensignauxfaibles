import Vue from 'vue'
import './plugins/vuetify'
import App from './App'
import router from './router'
import store from './store'

import axios from 'axios'

Vue.config.productionTip = false

Vue.prototype.$axios = axios.create(
  {
    headers: {
      'Content-Type': 'application/json'
    }
  }
)

Vue.prototype.$axios.interceptors.request.use(
  config => {
    config.baseURL = 'http://localhost:3000'
    if (store.sessionStore.state.token != null) config.headers['Authorization'] = 'Bearer ' + store.sessionStore.state.token
    return config
  }
)

Vue.prototype.$store = store.sessionStore
Vue.prototype.$localStore = store.localStore

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
