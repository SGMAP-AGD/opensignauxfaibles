<template>
  <v-card class="elevation-6" >
    <v-card-title>
      <v-list style="width: 100%">
        <v-list-tile >
          <v-list-tile-title>
            <v-combobox
              :items="files"
              label="Ajouter un fichier"
            ></v-combobox>
          </v-list-tile-title>
          <v-list-tile-content>
            <v-tooltip bottom>
            <v-btn icon slot="activator">
              <v-icon>fa-plus-square</v-icon>
            </v-btn>
            Ajouter
            </v-tooltip>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile
        v-for="file in currentFiles"
        :key="file"
        >
          <v-list-tile-content>
            {{ file }}
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
      
    </v-card-title>
  </v-card>
</template>

<script>
export default {
  props: ['type'],
  computed: {
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch () {
      if (this.$store.state.batches != null) {
        return this.$store.state.batches[this.currentBatchKey]
      }
    },
    currentFiles () {
      if (this.type != null && this.currentBatch != null) {
        return this.currentBatch.files[this.type.type]
      }
    },
    files () {
      return this.$store.state.files
    }
  }
}
</script>