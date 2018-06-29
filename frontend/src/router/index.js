import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import Data from '@/components/Data'
import Admin from '@/components/Admin'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/data',
      name: 'Data',
      component: Data
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin
    }
  ]
})
