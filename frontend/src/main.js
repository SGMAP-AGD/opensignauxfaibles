// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.

import Vue from 'vue'
import VueMaterial from 'vue-material'
import router from './router'
import App from './App'
import 'vue-material/dist/vue-material.min.css'
Vue.use(VueMaterial)

Vue.config.productionTip = false

/* eslint-disable no-new */

Vue.prototype.$api = 'http://localhost:3000/api/'

new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
