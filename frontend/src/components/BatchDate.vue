<template>
  <v-date-picker
  landscape
  locale="fr-FR"
  color="red darken-4"
  type="month"
  v-model="currentDate">
  </v-date-picker>
</template>

<script>
export default {
  props: ['date', 'param'],
  methods: {
    monthToISO (month) {
      return (new Date(Date.parse(month))).toISOString()
    }
  },
  computed: {
    currentBatch: {
      get () {
        if (this.$store.state.batches !== [] && this.currentBatchKey in this.$store.state.batches) {
          return this.$store.state.batches[this.currentBatchKey]
        } else {
          return {'params': {}}
        }
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch).then(r => this.$store.dispatch('checkEpoch'))
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
    currentDate: {
      get () {
        if (this.$store.state.batches != null) {
          var date = (this.currentBatch.params[this.param.prop] || '').substring(0, 7)
          date = (date < '1970-01') ? new Date().toISOString().substring(0, 7) : date
          return date
        }
      },
      set (month) {
        var batch = this.currentBatch
        batch.params[this.param.prop] = this.monthToISO(month)
        this.$store.dispatch('saveBatch', batch)
      }
    }
  }
}
</script>