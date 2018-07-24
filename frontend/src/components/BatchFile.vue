<template>
  <v-card>
    <v-toolbar  card dense>
      <v-toolbar-title>Lot {{ batch.id.key }} Fichiers</v-toolbar-title>
    </v-toolbar>      
    <v-card-title>
      <v-container grid-list-md text-xs-center>
        <v-layout row wrap>
          <v-flex xs6>
            Type<br/>
            <v-select
              v-model="type"
              :items="types"
              item-text="text"
              item-value="type"
              prepend-icon="fa-file-import"
              @change="updateFilter()"
              >
            </v-select>
          </v-flex>
          <v-flex xs6>
            Filtrer<br/>
            <v-text-field 
              v-model="filter"
              prepend-icon="fa-filter">
            </v-text-field>
          </v-flex>

        </v-layout>
      </v-container>
    </v-card-title>
    <v-card-text>
      <v-container grid-list-md text-xs-center>
        <v-layout row wrap>

          <v-flex xs6>
            <h2>Fichiers attachés</h2>
            <ul>
              <li
                v-for="file in batch.files[type]"
                v-bind:key="file"
                no-action>
              <h4>{{ file }}</h4>
              </li>
            </ul>
          </v-flex>
          <v-flex xs6>
            <v-checkbox
              v-for="file in files.filter(f => f.match(this.filter) && this.type != null)"
              v-bind:key="file"
              v-model="batch.files[type]"
              :disabled="type == null"
              :label="file"
              :value="file"
              @change="alterBatch()">
            </v-checkbox>
          </v-flex>
        </v-layout>
      </v-container>
    </v-card-text>
    <v-card-actions>
      <v-btn 
        flat
        v-on:click="close()">retour</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
  export default {
    name: 'BatchFile',
    methods: {
      test () {
        alert(null)
      },
      alterBatch () {
        this.batch.altered = true
      },
      close () {
        this.batch.dialog = false
      },
      updateFilter () {
        if (this.batch.files[this.type] == null) {
          this.batch.files[this.type] = []
        }
        this.filter = this.types.filter(i => (i.type === this.type))[0].filter
      }
    },
    props: {
      batch: {
        type: Object,
        default: () => ({})
      },
      types: {
        type: Array,
        default: () => ({})
      },
      files: {
        type: Array,
        default: () => ({})
      }
    },
    data () {
      return {
        filter: 'cotisation',
        candidates: [],
        selected: [],
        type: null
      }
    }
  }
</script>
