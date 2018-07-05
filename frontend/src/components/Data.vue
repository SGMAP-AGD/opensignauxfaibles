<template>
<div>
  <v-container grid-list-md text-xs-center>
    <v-layout>
      <v-flex key="datasource" xs3>
        <v-card>
          <v-card-title>
            <div>
              <span>Attachement fichiers</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
            v-model="attachBatchID"
            :items="batches"
            item-text="id.key"
            item-value="id.key"
            label="Batch"
            required
            ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-combobox
            :items="files"
            v-model="attachFilesID"
            label="Fichier"
            multiple
            chips
            ></v-combobox>
          </v-card-actions>
          <v-card-actions>
            <v-select
              :items="types"
              v-model="attachTypeID"
              label="Type"
            ></v-select>  
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="attachFiles()">Associer</v-btn>
          </v-card-actions>
        </v-card>
        <v-card>
          <v-card-title>
            <div>
              <span>Sources de données</span><br>
            </div>
            <v-card-actions>
              <v-text-field
                name="batchID"
                v-model="createBatchID"
                label="YYMM"
                value=""
                single-line
              ></v-text-field>  
            </v-card-actions>
            <v-card-actions>
              <v-btn @click="createBatch()">Créer lot</v-btn>
            </v-card-actions>
            <v-card-actions>
              <v-btn @click="listBatch()">Liste lot</v-btn>
            </v-card-actions>
          </v-card-title>
        </v-card>
      </v-flex>
      <v-flex key="import" xs6>
        <v-expansion-panel>
          <v-expansion-panel-content
            v-for="b in batches"
            :key="b.id.key"
          >
            <div slot="header">Lot {{b.id.key}} - {{ Object.keys(b.files).reduce((m,n) => { return m + b.files[n].length }, 0) }} fichiers</div>
            <v-card>
              <Batch v-bind:key="b.id.key" :batch="b"></Batch>
            </v-card>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-flex>
      <v-flex key="reduce" xs3>
        <v-card>
          <v-card-title>
            <div>
              <span>Import</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
            v-model="importBatchID"
            :items="batches"
            item-text="id.key"
            item-value="id.key"
            label="Batch"
            required
            ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="importData()">Import</v-btn>
          </v-card-actions>
        </v-card>
        <v-card>
          <v-card-title>
            <div>
              <span>Compactage</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-btn @click="compactEtablissement()">Etablissement</v-btn>
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="compactEntreprise()">Entreprise</v-btn>
          </v-card-actions>
        </v-card>
        <v-card>
          <v-card-title>
            <div>
              <span>Calcul des variables</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
              v-model="reduceFeatureID"
              :items="features"
              label="Algorithme"
              required
            ></v-select>
            <v-select
              v-model="reduceBatchID"
              :items="batches"
              item-text="id.key"
              item-value="id.key"
              label="Batch"
              required
            ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="reduce()">Reduce !</v-btn>
          </v-card-actions>
        </v-card>
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
    importData () {
      var self = this
      axios.get(this.$api + '/import/' + self.importBatchID)
    },
    createBatch () {
      var self = this
      axios.put(this.$api + '/admin/batch/' + this.createBatchID)
      .then(response => self.refresh())
    },
    refresh () {
      var self = this
      axios.get(this.$api + '/admin/batch').then(response => {
        self.batches = response.data
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
      name: 'App',
      createBatchID: '',
      features: [],
      importBatchID: '',
      reduceBatchID: '',
      reduceFeatureID: '',
      files: [],
      fileID: '',
      types: [],
      attachTypeID: '',
      batches: [],
      attachBatchID: '',
      attachFilesID: [],
      viewFiles: {}
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
