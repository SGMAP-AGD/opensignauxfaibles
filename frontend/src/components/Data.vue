<template>
<div>

  <v-container grid-list-md text-xs-center>
    <v-layout>
      <v-flex key="datasource" xs3>
        <v-card>
          <v-card-title>
            <div>
              <span>Sources de données</span><br>
            </div>
            <v-card-actions>
              <v-text-field
                name="batchID"
                v-model="batchID"
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

        <v-card>
          <v-card-title>
            <div>
              <span>Attachement fichiers</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
            v-model="selected_batch"
            :items="batch"
            label="Batch"
            required
            ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-select
              :items="files"
              v-model="file"
              label="Fichier"
          ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-select
              :items="types"
              v-model="type"
              label="Type"
            ></v-select>  
          </v-card-actions>
          
          <v-card-actions>
            <v-btn @click="attachFile()">Associer</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
      <v-flex key="import" xs3>
        <v-card>
          <v-card-title>
            <div>
              <span>Import</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
            v-model="selected_batch"
            :items="batch"
            label="Batch"
            required
            ></v-select>
          </v-card-actions>
          <v-card-actions>
            <v-btn @click="importData()">Import</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
      <v-flex key="compact" xs3>
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
      </v-flex>
      <v-flex key="reduce" xs3>
        <v-card>
          <v-card-title>
            <div>
              <span>Calcul des variables</span><br>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-select
              v-model="selected_algo"
              :items="algo"
              label="Algorithme"
              required
            ></v-select>
            <v-select
              v-model="selected_batch"
              :items="batch"
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
import axios from 'axios'
export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    },
    reduce () {
      axios.get('http://localhost:3000/api/reduce/' + this.selected_algo + '/' + this.selected_batch).then(response => alert('coucou'))
    },
    compactEtablissement () {
      axios.get('http://localhost:3000/api/compact/etablissement').then()
    },
    compactEntreprise () {
      axios.get('http://localhost:3000/api/compact/entreprise').then()
    },
    importData () {
      alert('Pas encore implémenté')
    },
    createBatch () {
      axios.put(this.$api + '/admin/batch/' + this.batchID)
      .then(response => alert(JSON.stringify(response.data, null, 2)))
    },
    listBatch () {
      var self = this
      axios.get(this.$api + '/admin/batch').then(response => { self.batch = response.data.map(batch => batch.id.key) })
    },
    listFiles () {
      var self = this
      axios.get(this.$api + '/admin/files').then(response => { self.files = response.data })
    },
    listTypes () {
      var self = this
      axios.get(this.$api + '/admin/types').then(response => { self.types = response.data })
    },
    attachFile () {
      var self = this
      axios.post(this.$api + '/admin/attach',
        {
          'file': self.file,
          'type': self.type,
          'batch': self.selected_batch
        })
    }
  },
  mounted () {
    this.listBatch()
    this.listFiles()
    this.listTypes()
  },
  data () {
    return {
      batchID: '',
      importBatch: '',
      batch: [],
      algo: [],
      name: 'App',
      files: [],
      file: '',
      selected_batch: '',
      selected_algo: '',
      types: [],
      type: ''
    }
  }
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
