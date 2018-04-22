// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.

import Vue from 'vue'
import Vuex from 'vuex'
import VueMaterial from 'vue-material'
import App from './App'
import 'vue-material/dist/vue-material.min.css'
import router from './router'
import store from '@/store/store'

Vue.use(Vuex)
Vue.use(VueMaterial)

Vue.config.productionTip = false

/* eslint-disable no-new */

Vue.prototype.$api = 'http://localhost:3000/api'
Vue.prototype.$generatePeriodSerie = function (date_debut, date_fin) {
  var date_next = new Date(date_debut.getTime())
  var serie = []
  while (date_next.getTime() < date_fin.getTime()) {
    serie.push(new Date(date_next.getTime()))
    date_next.setUTCMonth(date_next.getUTCMonth() + 1)
  }return serie
}

new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>',
  beforeCreate: function () {
  }
})
