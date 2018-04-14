<template>
<div>
    <md-card>
      <md-card-header>
        <div class="md-title">Visualisation de débit</div>
        <div class="md-subhead">by Signaux-Faibles™</div>
        
      </md-card-header>
      <md-card-content>
        <md-field>
          <label>Initial Value</label>
          <md-input v-model='siret'/>
        </md-field>

        <md-button class='md-raised md-primary' v-on:click='fillTable()'>Consulter</md-button>
        
        <md-chips v-model="sirets" md-placeholder="Add siret..."></md-chips>

        <vue-plotly v-for='siret in sirets' v-bind:key='siret' :data='data' :layout='layout' :options='options'/>
      </md-card-content>
    </md-card>
    
</div>
</template>

<script>
import VuePlotly from '@statnett/vue-plotly'
import axios from 'axios'

export default {
  name: 'DataDebit',
  components: {
    VuePlotly
  },
  data () {
    return {
      sirets: [],
      data: [],
      siret: '',
      layout: {barmode: 'stack'},
      options: {}
    }
  },
  methods: {
    fillTable: function() {
      axios.get(`http://localhost:3000/api/v1/data/debit/` + this.siret)
        .then(response => {this.data = response.data[0].value})
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
