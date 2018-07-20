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
  <v-container  grid-list-md text-xs-center>
    <v-layout justify-start row wrap>
      <v-flex 
       xs4
       >
       <Batch :batch="newBatch" newBatch/>
      </v-flex>
      <v-flex 
       xs4
       v-for="batch in batches"
       :key="batch.id.key"
       >
        <Batch :batch="batch" v-if="(viewReadonly)?true:!batch.readonly"/>
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
    reduce () {
      axios.get(this.$api + '/reduce/' + this.reduceFeatureID + '/' + this.reduceBatchID).then(response => alert('coucou'))
    },
    compactEtablissement () {
      axios.get(this.$api + '/compact/etablissement').then()
    },
    compactEntreprise () {
      axios.get(this.$api + '/compact/entreprise').then()
    },
    refresh () {
      var self = this
      axios.get(this.$api + '/admin/batch').then(response => {
        this.newBatch = {
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
        self.batches = response.data.map(b => {
          b.altered = false
          b.params.date_debut = (b.params.date_debut === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_debut.substring(0, 7)
          b.params.date_fin = (b.params.date_fin === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_fin.substring(0, 7)
          b.params.date_fin_effectif = (b.params.date_fin_effectif === '0001-01-01T00:00:00Z') ? (new Date()).toISOString().substring(0, 7) : b.params.date_fin_effectif.substring(0, 7)
          return b
        })
      })
      axios.get(this.$api + '/admin/files').then(response => {
        this.files = response.data.map(file => {
          return {
            'text': file.split('/').reverse().slice(0, 3).reverse().join('/'),
            'value': file
          }
        })
      })
      axios.get(this.$api + '/admin/types').then(response => { self.types = response.data.sort() })
      axios.get(this.$api + '/admin/features').then(response => { self.features = response.data })
    },
    attachFiles () {
      this.attachFilesRecur(this.attachFilesID)
    },
    attachFilesRecur (files) {
      var self = this
      if (files.length > 0) {
        axios.post(this.$api + '/admin/attach',
          {
            'file': files[0].value,
            'type': self.attachTypeID,
            'batch': self.attachBatchID
          }).then(response => { this.attachFilesRecur(files.slice(1)) })
      } else {
        this.refresh()
      }
    }
  },
  mounted () {
    this.refresh()
  },
  data () {
    return {
      'name': 'App',
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
<!-- Add "scoped" attribute to limit CSS to this component only -->
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
