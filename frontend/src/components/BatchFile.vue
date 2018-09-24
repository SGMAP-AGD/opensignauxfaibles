<template>

  <v-card class="elevation-6" >
    <v-card-title class="header">
      <v-toolbar
        class="elevation-3"
        color="red darken-4"
        dark
        card
      >
      <v-toolbar-title>{{ currentType.text }}</v-toolbar-title>
    </v-toolbar>
    </v-card-title>
    <v-card-text>
      <v-flex xs12>
        <v-combobox 
          multiple
          deletableChips
          chips
          flat
          v-model="addFiles"
          :items="files"
          item-text="filename"
          item-value="name"
          label="Ajouter un fichier"
          ripple
          
          append-outer-icon="fa-plus-square"
          @click:append-outer="add"
        >
          <template
            slot="item"
            slot-scope="{index, item, parent}"
          >
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
          </template>
        </v-combobox>
      </v-flex>
      <v-list style="width: 100%">
        <v-list-tile style="width: 100%">
          <v-list-tile-content>
            <span class="strong">{{ currentFiles.length }} fichier{{ (currentFiles.length>0?"s":"") }}</span>
          </v-list-tile-content>
          <v-list-tile-action>
            <span class="light">{{ formatBytes(currentFiles.reduce((m,f) => m+f.size, 0)) }}</span>
          </v-list-tile-action> 
        </v-list-tile>

        <v-divider></v-divider>
        <v-list-tile style="width: 100%"
        v-for="file in currentFiles"
        :key="file.name"
        ripple
        @click="removeMark(file.name)"
        :class="removeFiles.includes(file.name)?'todelete':'tokeep'"
        >
          <v-list-tile-content>
            <span class="strong">{{ file.filename }}</span>
            <span class="light">{{ file.pathname }}</span>
          </v-list-tile-content>
          <v-list-tile-action>
            <span class="light">{{ file.psize }}</span>
            <span class="light">{{ file.pdate }}</span>
          </v-list-tile-action>
        </v-list-tile>
      </v-list>
    </v-card-text>
    <v-card-action v-if="removeFiles.length > 0">
      <v-btn 
      block 
      flat 
      v-if="removeFiles.length > 0"
      @click="remove()"
      >
      <v-icon >fa-trash</v-icon>
      </v-btn>
    </v-card-action>
  </v-card>
</template>

<script>
export default {
  props: ['type'],
  data () {
    return {
      addFiles: [],
      removeFiles: []
    }
  },
  methods: {
    comboFilter (item, queryText, itemText) {
      return item.name.toLowerCase().includes(queryText.toLowerCase())
    },
    formatBytes (a, b) {
      if (a === 0) return '0 Bytes'
      var c = 1024
      var d = b || 2
      var e = ['Octets', 'Kio', 'Mio', 'Gio', 'Tio', 'Pio', 'Eio', 'Zio', 'Yio']
      var f = Math.floor(Math.log(a) / Math.log(c))
      return parseFloat((a / Math.pow(c, f)).toFixed(d)) + ' ' + e[f]
    },
    formatDate (d) {
      return d.getFullYear() + '/' + d.getMonth() + '/' + d.getDate()
    },
    add () {
      var batch = this.currentBatch
      this.addFiles.forEach(f => {
        if (!(batch.files[this.type] || []).includes(f.name)) {
          batch.files[this.type] = (batch.files[this.type] || []).concat(f.name)
        }
      })
      this.addFiles = []
      this.currentBatch = batch
    },
    fileDetail (file) {
      var candidates = (this.files || []).filter(f => f.name === file)
      if (candidates.length === 1) {
        return candidates[0]
      } else {
        var arrayFile = file.split('/')
        var lengthArrayFile = arrayFile.length
        return {
          name: file,
          pathname: arrayFile.slice(0, lengthArrayFile - 1).join('/') + '/',
          filename: arrayFile[lengthArrayFile - 1],
          psize: 'n/a',
          pdate: 'n/a',
          date: 'n/a'
        }
      }
    },
    removeMark (file) {
      if (this.removeFiles.includes(file)) {
        this.removeFiles = this.removeFiles.filter(f => f !== file)
      } else {
        this.removeFiles.push(file)
      }
    },
    remove () {
      var batch = this.currentBatch
      batch.files[this.type] = batch.files[this.type].filter(f => !this.removeFiles.includes(f))
      this.removeFiles = []
      this.currentBatch = batch
    }
  },
  computed: {
    currentType () {
      return this.$store.state.types.filter(t => t.type === this.type)[0]
    },
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch: {
      get () {
        if (this.$store.state.batches !== []) {
          return this.$store.state.batches[this.currentBatchKey]
        } else {
          return null
        }
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch).then(r => this.$store.dispatch('checkEpoch'))
      }
    },
    currentFiles () {
      if (this.$store.state.batches !== []) {
        return (this.currentBatch.files[this.type] || []).map(f => this.fileDetail(f))
      }
    },
    files () {
      var files = (this.$store.state.files || []).map(f => {
        var arrayFile = f.name.split('/')
        var lengthArrayFile = arrayFile.length
        f = {
          name: f.name,
          size: f.size,
          pathname: arrayFile.slice(0, lengthArrayFile - 1).join('/') + '/',
          filename: arrayFile[lengthArrayFile - 1],
          psize: this.formatBytes(f.size),
          pdate: f.date.substring(0, 10),
          date: new Date(f.date)
        }
        return f
      })
      return files
    },
    filterFiles () {
      return this.files.filter(f => {
        return !(this.currentFiles.includes(f))
      })
    }
  }
}
</script>

<style scoped>
.todelete {
  background-color: #FCE4EC;
}
.tokeep {
  background-color: white;
}
.light {
  color: #888;
  font-size: 10px;
}
.string {
  color: #000;
  font-size: 11px;
}
</style>