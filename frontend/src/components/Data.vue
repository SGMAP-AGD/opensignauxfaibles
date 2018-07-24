<template>
<div>
<v-container  grid-list-md text-xs-center>
    <v-layout justify-start row wrap>
      
      <v-flex 
      xs12
      >
        <v-toolbar>
          <v-switch label="Lots verrouillés" v-model="viewReadonly"></v-switch>
          <v-spacer></v-spacer><v-btn>compacter</v-btn>
        </v-toolbar>
      </v-flex>
    </v-layout>
  </v-container>
  <v-container grid-list-md text-xs-center>
    <v-layout justify-start row wrap >
        <v-flex 
        xs4
        >
        
          <Batch :batch="newBatch" :types="types" :files="files" newBatch/>
        
        </v-flex>
        <v-flex 
        xs4
        v-for="batch in batches"
        :key="batch.id.key"
        >
        <v-fade-transition>
          <Batch :batch="batch" :types="types" :files="files" v-if="(viewReadonly)?true:!batch.readonly"/>
        </v-fade-transition>
        </v-flex>
    </v-layout>
  </v-container>
</div>
</template>

<script>
import Batch from '@/components/Batch'
import axios from 'axios'

export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    },
    refresh () {
      var self = this
      axios.get(this.$api + '/admin/types').then(response => {
        self.types = response.data
        axios.get(this.$api + '/admin/batch').then(response => {
          this.newBatch = {
            'id': {
              'key': '',
              'type': 'batch'
            },

            'files': self.types.reduce((a, h) => {
              a[h.type] = []
              return a
            }, {}),
            'log': null,
            'readonly': false,
            'params': {
              'date_debut': (new Date()).toISOString().substring(0, 7),
              'date_fin': (new Date()).toISOString().substring(0, 7),
              'date_fin_effectif': (new Date()).toISOString().substring(0, 7)
            }
          }
          self.batches = response.data.map(b => {
            b.altered = false
            b.params.date_debut = (b.params.date_debut === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_debut.substring(0, 7)
            b.params.date_fin = (b.params.date_fin === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_fin.substring(0, 7)
            b.params.date_fin_effectif = (b.params.date_fin_effectif === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_fin_effectif.substring(0, 7)
            b.dialog = false
            b.files = self.types.reduce((a, h) => {
              a[h.type] = (b.files[h.type] || [])
              return a
            }, {})
            return b
          })
        })
      })

      axios.get(this.$api + '/admin/types').then(response => { self.types = response.data.sort() })
      axios.get(this.$api + '/admin/features').then(response => { self.features = response.data })

      axios.get(this.$api + '/admin/files').then(response => {
        self.files = response.data
      })
    }
  },
  mounted () {
    this.refresh()
  },
  data () {
    return {
      'types': [],
      'files': [],
      'viewReadonly': false,
      'batches': [],
      'newBatch': {
        'id': {
          'key': '',
          'type': 'batch'
        },
        'files': {},
        'log': null,
        'readonly': false,
        'params': {
          'date_debut': (new Date()).toISOString().substring(0, 7),
          'date_fin': (new Date()).toISOString().substring(0, 7),
          'date_fin_effectif': (new Date()).toISOString().substring(0, 7)
        }
      }
    }
  },
  components: { Batch }
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
