<template>
  <div id="app">
    <md-card>
      <md-card-header>
        <div class="md-title">Front-end Signaux-Faibles</div>
        <div class="md-subhead">You can't test</div>
      </md-card-header>

      <md-card-content>
        <md-button class="md-raised md-primary" v-on:click="lookFiles()">Look for files</md-button>
        <md-button class="md-raised md-primary" v-on:click="clearFiles()">Clear List</md-button>
        <br />

        <md-table>
            <md-table-row>
              <md-table-head>Filename</md-table-head>
              <md-table-head>Path</md-table-head>
              <md-table-head>Date</md-table-head>
            </md-table-row>
            <md-table-row v-for="file in files" :key="file.id" >
              <md-table-cell>{{file.name}}</md-table-cell>
              <md-table-cell>{{file.path}}</md-table-cell>
              <md-table-cell>{{file.date}}</md-table-cell>
            </md-table-row>
        </md-table>
      </md-card-content>
    </md-card>
  </div>
</template>

<script>

export default {
  name: 'app',
  data: function () {
    return {
      files: []
    }
  },
  methods: {
    lookFiles: function () {
      this.axios.get('http://localhost:3000/api/v1/listFiles').then((response) => {
        response.data.map(file => { this.files.push(file) })
      })
    },
    clearFiles: function () {
      this.files.splice(0, this.files.length)
    }
  }
}
</script>
