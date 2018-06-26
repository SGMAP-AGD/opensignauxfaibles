import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import DataImport from '@/components/DataImport'
import DataBatch from '@/components/DataBatch'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/data/import',
      name: 'DataImport',
      component: DataImport
    },
    {
      path: '/data/batch',
      name: 'DataBatch',
      component: DataBatch
    }
  ]
})
