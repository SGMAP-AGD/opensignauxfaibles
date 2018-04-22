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
  <div id="app">
    <chart-debit :siret="siret"></chart-debit>
</div>
      </md-card-content>
    </md-card>

</div>
</template>

<script>
import ChartDebit from '@/components/ChartDebit'
export default {
  name: 'DataDebit',
  components: { ChartDebit },
  data () {
    return {
      siret: '',
      data: []
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
