<template>
<div>
  <v-card>
    <v-card-title
      class="headline grey lighten-2"
      primary-title
    >
      <v-container>
        <v-layout>
          <v-flex xs6>
            {{ batch.id.key }} - Fichiers
          </v-flex>
          <v-flex xs6>
            <v-select
            v-model="type"
            :items="types"
            item-text="text"
            item-value="type">
            </v-select>
          </v-flex>
        </v-layout>
      </v-container>
    </v-card-title>

    <v-card-text>
      <h2>{{ type }}</h2>
      <ul>
          <li                              
          v-for="file in batch.files[type]"
          v-bind:key="file"
          no-action>
          <h4>{{ file }}</h4>
        </li>
      </ul>
    </v-card-text>
  </v-card>


  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'BatchFile',
  methods: {
    test () {
      alert(null)
    }
  },
  mounted () {
    var self = this
    axios.get(this.$api + '/admin/files').then(response => {
      self.files = response.data
    })
    axios.get(this.$api + '/admin/types').then(response => {
      self.types = response.data
    })
  },
  props: {
    batch: {
      type: Object,
      default: () => ({})
    }
  },
  data () {
    return {
      files: [],
      types: [],
      type: null
    }
  }
}
</script>
