<template>
<div>
  <v-container grid-list-md text-xs-center>
    <v-layout row wrap>
      <v-flex xs3>
        <v-navigation-drawer
          permanent
          float
          raised
          style="z-index: 1"
          >
            <v-list dense class="pt-0">
              <v-list-group>
                <v-list-tile slot="activator">
                  <v-list-tile-content class="title">
                    Paramètres
                  </v-list-tile-content>
                </v-list-tile>
                <v-list-tile
                v-for="param in parameters"
                :key="param.key"
                ripple
                @click="setCurrentType(param.key)">
                  <v-list-tile-content
                  :class="(param.key===currentType) ? 'selected': null"
                  >
                    <v-list-tile-title>{{ param.text }}</v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list-group>
              <v-list-group>
                <v-list-tile slot="activator">
                  <v-list-tile-title class="title">
                    Fichiers
                  </v-list-tile-title>
                </v-list-tile>

                <v-divider></v-divider>
                <v-list-tile
                v-for="type in types"
                :key="type.text"
                @click="setCurrentType(type)"
                >
                  <v-list-tile-content
                  :class="(type==currentType) ? 'selected': null"
                  >
                    <v-list-tile-title>{{ type.text }}</v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list-group>
            </v-list>
          </v-navigation-drawer>
      </v-flex>
      <v-flex xs9>
        <BatchDate 
        :date="currentType"
        :param="parameters.filter(p => p.key === currentType)[0]"
        v-if="parameters.map(p => p.key).includes(currentType)"
        />
        <BatchFile 
        :type="currentType"
        v-if="types.includes(currentType)"
        />
      </v-flex>
    </v-layout>
  </v-container>
    </div>
</template>

<script>
import BatchFile from '@/components/BatchFile'
import BatchDate from '@/components/BatchDate'

export default {
  props: ['batchKey'],
  data () {
    return {
      currentType: null,
      parameters: [
        {text: 'Date de début', key: 'dateDebut', prop: 'date_debut'},
        {text: 'Date de fin', key: 'dateFin', prop: 'date_fin'},
        {text: 'Date de fin effectifs', key: 'dateFinEffectif', prop: 'date_fin_effectif'}
      ]
    }
  },
  methods: {
    setCurrentType (type) {
      this.currentType = type
    }
  },
  computed: {
    currentBatch () {
      return this.$store.state.batches.filter(b => b.id.key === this.batchKey)
    },
    types () {
      return this.$store.state.types
    },
    features () {
      return this.$store.state.features
    },
    files () {
      return this.$store.state.files
    }
  },
  components: { BatchFile, BatchDate }
}
</script>

<style>
  .selected {
    color: blue;
    font-size: 14px;
  }
</style>