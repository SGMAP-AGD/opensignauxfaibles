// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import App from './App'
import router from './router'
import VueMaterial from 'vue-material'
import 'vue-material/dist/theme/default-dark.css' // This line here
import 'vue-material/dist/vue-material.css'

Vue.use(VueAxios, axios)
Vue.use(VueMaterial)
Vue.config.productionTip = false

Vue.prototype.$endpoint = 'localhost:3000'
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
