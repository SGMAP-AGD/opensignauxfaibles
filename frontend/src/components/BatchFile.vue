<template>

  <v-card class="elevation-6" >
    <v-card-title>
      <v-flex xs11>
        <v-combobox style="width: 90%"
          multiple
          deletableChips
          chips
          outline
          v-model="addFiles"
          :items="files"
          item-text="filename"
          label="Ajouter un fichier"
        >
          <template
            slot="item"
            slot-scope="{index, item, parent}"
          >
            <v-list dense style="width: 100%">
              <v-list-tile style="width: 100%">
                <v-list-tile-content>
                  <span class="strong">{{ item.filename }}</span>
                  <span class="light">{{ item.pathname }}</span>
                </v-list-tile-content>
                <v-list-tile-action>
                  <span class="light">{{ item.psize }}</span>
                  <span class="light">{{ item.pdate }}</span>
                </v-list-tile-action>
              </v-list-tile>
            </v-list>
          </template>
        </v-combobox>
      </v-flex>
      <v-flex xs1>
        <div class="text-xs-center">
          <v-tooltip bottom class="text-xs-center">
            <v-btn  class="text-xs-center"
            large 
            fab
            slot="activator"
            @click="add()">
              <v-icon large>fa-plus-square</v-icon>
            </v-btn>
            Ajouter
          </v-tooltip>
        </div>
      </v-flex>
      <v-list style="width: 100%">
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
  methods: {
    formatBytes (a, b) {
      if (a === 0) return '0 Bytes'
      var c = 1024
      var d = b || 2
      var e = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
      var f = Math.floor(Math.log(a) / Math.log(c))
      return parseFloat((a / Math.pow(c, f)).toFixed(d)) + ' ' + e[f]
    },
    formatDate (d) {
      return d.getFullYear() + '/' + d.getMonth() + '/' + d.getDate()
    },
    add () {
      var batch = this.currentBatch
      this.addFiles.forEach(f => {
        batch.files[this.type] = (batch.files[this.type] || []).concat(f.name)
      })
      this.addFiles = []
      this.currentBatch = batch
    }
  },
  computed: {
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch: {
      get () {
        if (this.$store.state.batches != null) return this.$store.state.batches[this.currentBatchKey]
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch)
      }
    },
    currentFiles () {
      if (this.type != null && this.currentBatch != null) {
        return this.currentBatch.files[this.type]
      }
    },
    files () {
      return this.$store.state.files.map(f => {
        var arrayFile = f.name.split('/')
        var lengthArrayFile = arrayFile.length
        f = {
          name: f.name,
          pathname: arrayFile.slice(0, lengthArrayFile - 1).join('/') + '/',
          filename: arrayFile[lengthArrayFile - 1],
          psize: this.formatBytes(f.size),
          pdate: f.date.substring(0, 10),
          date: new Date(f.date)
        }
        return f
      }).sort((a, b) => (a.date.getTime() < b.date.getTime()) ? 1 : -1)
    }
  },
  data () {
    return {
      addFiles: []
    }
  }
}
</script>

<style scoped>
.light {
  color: #888;
  font-size: 10px;
}
.string {
  color: #000;
  font-size: 11px;
}
</style>