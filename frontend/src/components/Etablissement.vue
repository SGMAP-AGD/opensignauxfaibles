<template>
  <div>
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    style="z-index: 1"
    >
      
      
      <v-img src="/static/logo.png"> </v-img>
    </v-navigation-drawer>
    <div class="widget">
      Etablissement {{ param }}
    </div>
  </div>
</template>

<script>
export default {
  props: ['param'],
  name: 'Etablissement'
}
</script>

<style scoped>
.echarts {
  width: 400px
}
.widget {
  position: absolute;
  left: 320px;
  top: 20px; 
  right: 20px;
}
</style>