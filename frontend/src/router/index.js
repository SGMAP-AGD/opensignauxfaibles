import Vue from 'vue'
import Router from 'vue-router'
import DataDebit from '@/components/DataDebit'
// import Plotly from '@/components/Plotly'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'DataDebit',
      component: DataDebit
    }
  ]
})
