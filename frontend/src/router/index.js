import Vue from 'vue'
import Router from 'vue-router'
import DataDebit from '@/components/DataDebit'
import Region from '@/components/Region'
import Landing from '@/components/Landing'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: Landing
    },
    {
      path: '/region',
      name: 'Region',
      component: Region
    },
    {
      path: '/data/debit',
      name: 'DataDebit',
      component: DataDebit
    }

  ]
})
