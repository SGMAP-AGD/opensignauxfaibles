<template>
  <div id="app">
    <div class="md-title">Open Signaux-Faibles™</div>
    <div class="md-subhead">Un défi Entrepreneur d'Intérêt Général</div>
    <Login id='login'/>
    <router-view/>
  </div>
</template>

<script>
import Login from '@/components/Login'

export default {
  name: 'App',
  components: { Login }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 10px;
}

</style>
