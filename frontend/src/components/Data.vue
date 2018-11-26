<template>
  <div>
<div>
  <v-toolbar height="35px" class="toolbar" color="#ffffff"  app>
    <v-icon 
    @click="drawer=!drawer"
    class="fa-rotate-180"
    v-if="!drawer"
    color="secondary"
    key="toolbar"
    >
    mdi-backburger
     </v-icon>
    <div style="width: 100%; text-align: center;"  class="titre">
      Données
    </div>
    <v-spacer></v-spacer>
    <v-icon color="secondary"  @click="rightDrawer=!rightDrawer">fa-database</v-icon>
  </v-toolbar>
  <div style="width:100%">
    <Batch
    batchKey="1802"/>
  </div>
    <v-tabs
      v-model="currentBatchKey"
      color="red darken-4"
      dark
      fixed
      slider-color="red accent-2"
      show-arrows
      next-icon="fa-arrow-alt-circle-right"
      prev-icon="fa-arrow-alt-circle-left"
    >
      <v-tab
        v-for="batch in batches"
        :key="batch"
        ripple
        dark
      >
        {{ batch.substring(2,4) }}/20{{ batch.substring(0,2) }}
      </v-tab>
      <v-tab-item
      dark
      v-for="(batch, rank) in batches"
      :key="rank"
      >
          <Batch 
          :batchKey="batch"/>
      </v-tab-item>
    </v-tabs> 
  </div>
  </div>
</template>

<script>
import Batch from '@/components/Batch'
export default {
  mounted () {
    this.$store.commit('updateBatches')
  },
  methods: {
    setCurrentBatchKey (batchKey) {
      this.currentBatchKey = batchKey
    }
  },
  computed: {
    drawer: {
      get () {
        return this.$store.state.appDrawer
      },
      set (val) {
        this.$store.dispatch('setDrawer', val)
      }
    },
    rightDrawer: {
      get () {
        return this.$store.state.rightDrawer
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
    currentBatchKey: {
      get () {
        return this.$store.state.currentBatchKey
      },
      set (value) {
        this.$store.commit('setCurrentBatchKey', value)
      }
    },
    batches () {
      return this.$store.state.batches.map(batch => batch.id.key)
    }
  },
  components: { Batch },
  name: 'Data'
}
</script>

<style scoped>
  div.titre {
    color: #8e0000;
    font-family: 'Signika', sans-serif;
    font-weight: 500;
    color: primary;
    font-size: 18px;
  }
</style>