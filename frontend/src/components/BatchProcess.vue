<template>
<v-card class="elevation-6">
  <v-card-text>
    
  </v-card-text>
  <v-card-media
    contain
    height="80px"
    :src="process.img"
  >
    
  </v-card-media>
  <v-card-text>
    {{ process.description }}
  </v-card-text>
  <v-card-actions>
    <v-btn
    :disabled="dbstatus != null"
    @click.stop="dialog = true"
    block>LANCER</v-btn>
    
    <v-dialog
      v-model="dialog"
      max-width="290"
    >
      <v-card>

        <v-card-title class="title">
          Êtes vous sûr ?
        </v-card-title>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            color="error"
            flat="flat"
            @click="dialog = false"
          >
            RETOUR
          </v-btn>

          <v-btn
            color="green darken-1"
            flat="flat"
            @click="mock(); dialog = false"
          >
            LANCER
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card-actions>
</v-card>
</template>

<script>
export default {
  props: ['process'],
  data () {
    return {
      dialog: false
    }
  },
  computed: {
    dbstatus () {
      return this.$store.state.dbstatus
    }
  },
  methods: {
    mock () {
      this.$axios.get('http://localhost:3000/api/mock/compact')
    }
  }
}
</script>