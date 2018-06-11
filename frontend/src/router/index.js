import Vue from 'vue'
import Router from 'vue-router'
import DataDebit from '@/components/DataDebit'
import DataView from '@/components/DataView'
import Region from '@/components/Region'
import Landing from '@/components/Landing'
import Tasks from '@/components/Tasks'

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
      path: '/tasks',
      name: 'Tasks',
      component: Tasks
    },
    {
      path: '/data/debit',
      name: 'DataDebit',
      component: DataDebit
    },
    {
      path: '/data/view',
      name: 'DataView',
      component: DataView
    }

  ]
})
