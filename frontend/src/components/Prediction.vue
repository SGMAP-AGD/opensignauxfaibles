<template>
<div>
  <v-tabs
    v-model="currentBatchKey"
    color="indigo darken-4"
    dark
    slider-color="red accent-2"
    show-arrows
    next-icon="fa-arrow-alt-circle-right"
    prev-icon="fa-arrow-alt-circle-left"

  >
    <v-tab
      v-for="batch in batches"
      :key="batch"
      ripple
    >
      {{ batch.substring(2,4) }}/20{{ batch.substring(0,2)}}
    </v-tab>
    <v-tab-item
      v-for="batch in batches"
      :key="batch"
      lazy
    >
        <PredictionTable :batchKey="batch"/>
    </v-tab-item>
  </v-tabs> 
</div>
</template>

<script>
import PredictionTable from '@/components/PredictionTable'

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
  components: { PredictionTable },
  name: 'Prediction'
}
</script>