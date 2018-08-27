<template>
  <v-date-picker
  color="red darken-4"
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
          return this.currentBatch.params[this.param.prop].substring(0, 10)
        }
      },
      set (date) {

      }
    }
  }
}
</script>