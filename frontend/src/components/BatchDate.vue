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
  computed: {
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch () {
      if (this.$store.state.batches != null) {
        return this.$store.state.batches[this.currentBatchKey]
      }
    },
    currentDate: {
      get () {
        if (this.$store.state.batches != null) {
          var date = this.currentBatch.params[this.param.prop].substring(0, 7)
          date = (date < '1970-01') ? '1970-01' : date
          console.log(date)
          return date
        }
      },
      set (date) {

      }
    }
  }
}
</script>