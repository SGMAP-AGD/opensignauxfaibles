<template>
<div>
    <div>
    <md-table v-model="tasks" md-sort="name" md-sort-order="asc" md-card @sort="onSort">
      <md-table-toolbar>
        <h1 class="md-title">Visites d'entreprises</h1>
      </md-table-toolbar>

      <md-table-row slot="md-table-row" slot-scope="{ item }">
        <md-table-cell md-label="ID" md-numeric>{{ item.id }}</md-table-cell>
        <md-table-cell md-label="Entreprise" md-sort-by="title">{{ item.title }}</md-table-cell>
        <md-table-cell md-label="Description" md-sort-by="description">{{ item.description }}</md-table-cell>
      </md-table-row>
    </md-table>
  </div>
</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'tasks',
  data () {
    return {
      tasks: []
    }
  },
  mounted: function () {
    var self = this
    axios.get(this.$api + '/kanboard/get/tasks').then(function (response) {
      self.tasks = response.data
    })
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
