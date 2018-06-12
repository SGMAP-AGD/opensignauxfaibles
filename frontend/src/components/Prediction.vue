<template>
<div>
    <md-card>
      <md-card-header>
        <div class="md-title">Visualisation des données entreprises</div>

      </md-card-header>
      <md-card-content>

      <md-field>
        <label for="batch">Lot</label>
        <md-select v-model="selectedBatch" name="batch" id="batch">
          <md-option v-for="b in batch" v-bind:key="b" :value="b">{{b}}</md-option>
        </md-select>
      </md-field>

      <md-field>
        <label for="algo">Algorithme</label>
        <md-select v-model="selectedAlgo" name="algo" id="algo">
          <md-option v-for="a in algo" v-bind:key="a" :value="a">{{a}}</md-option>
        </md-select>
      </md-field>

      <md-button class="md-dense md-raised md-primary" v-on:click="get()">Consulter</md-button>
      </md-card-content>
    </md-card>
</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Prediction',
  data () {
    return {
      batch: [],
      algo: [],
      selectedBatch: '',
      selectedAlgo: '',
      prediction: []
    }
  },
  methods: {
    get () {
      var self = this
      axios.get(this.$api + '/data/prediction/' + this.selectedAlgo + '/' + this.selectedBatch).then(response => {
        self.prediction = response.data
      })
    }
  },
  mounted () {
    var self = this
    axios.get(this.$api + '/data/batch').then(response => {
      self.batch = response.data
    })
    axios.get(this.$api + '/data/algo').then(response => {
      self.algo = response.data
    })
  }
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
  display: block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
