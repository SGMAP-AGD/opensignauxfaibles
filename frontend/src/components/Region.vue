<template>
<div>
    <md-card>
      <md-card-header>
        <div class="md-title">Administration des région</div>
        Blablabla
      </md-card-header>
      <md-card-content>

      </md-card-content>
    </md-card>

</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Region',
  data () {
    return {
      region: []
    }
  },
  methods: {
    getAdminRegion: function () {
      axios.get(`http://localhost:3000/api/v1/admin/region/`)
        .then(response => { this.data = response.data[0].value })
    }
  },
  beforeMounts: function () { this.getAdminRegion() }
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
