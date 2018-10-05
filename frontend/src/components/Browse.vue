<template>
<div>

        
  <v-tabs
      v-model="active"
      color="indigo darken-4"
      dark
      slider-color="red accent-2"
    >
      <v-tab
        v-for="(tab, index) in tabs"
        :key="index"
        ripple
      >
       {{ tab.param }}
      </v-tab>
      <v-tab-item
        style="min-height: 500vh;"
        v-for="(tab,index) in tabs"
        :key="index"
      >
        <PredictionTable v-if="tab.type==='Prediction'" :param="tab.param"/>
        <Etablissement v-if="tab.type==='Etablissement'" :param="tab.param"/>
      </v-tab-item>
    </v-tabs>

</div>
</template>

<script>
import PredictionTable from '@/components/PredictionTable'
import Etablissement from '@/components/Etablissement'

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
    tabs: {
      set (tabs) {

      },
      get () {
        return this.$store.getters.getTabs
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
      return this.$store.state.batches.filter(b => b.readonly === true).map(batch => batch.id.key)
    }
  },
  components: { PredictionTable, Etablissement },
  name: 'Browse'
}
</script>

