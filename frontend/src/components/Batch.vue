<template>
<div class="container">
  <div class="fixed d-inline-block elevation-6">
    <v-navigation-drawer
    permanent
    
    style="z-index: 1"
    >
      <v-list dense class="pt-0">
        <v-list-group>
          <v-list-tile slot="activator" bgcolor="red">
            <v-list-tile-action>
              <v-icon>fa-cogs</v-icon>
            </v-list-tile-action>
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
            <v-list-tile-action>
              <v-icon>fa-copy</v-icon>
            </v-list-tile-action> 
            <v-list-tile-title class="title">
              Fichiers
            </v-list-tile-title>
          </v-list-tile>

          <v-divider></v-divider>
          <v-list-tile
          v-for="type in types"
          :key="type.text"
          @click="setCurrentType(type.type)"
          >
            <v-list-tile-content
            :class="(type.type==currentType) ? 'selected': null"
            >
              <v-list-tile-title>{{ type.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
        <v-list-group>
          <v-list-tile slot="activator">
            <v-list-tile-action>
              <v-icon>fa-microchip</v-icon>
            </v-list-tile-action> 
            <v-list-tile-title class="title">
              Traitements
            </v-list-tile-title>
          </v-list-tile>
          <v-divider></v-divider>
          <v-list-tile
          v-for="process in processes"
          :key="process.key"
          @click="setCurrentType(process.key)"
          >
            <v-list-tile-content
            :class="(process.key==currentType) ? 'selected': null"
            >
              <v-list-tile-title>{{ process.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
      </v-list>
    </v-navigation-drawer>
    </div>
    <div class="flex-item">
      <v-container grid-list-xs text-xs-center>
        <v-layout row wrap justify-start>
          <v-flex xs12 >
            <BatchDate 
            class="d-inline-block elevation-6"

            :date="currentType"
            :param="parameters.filter(p => p.key === currentType)[0]"
            v-if="parameters.map(p => p.key).includes(currentType)"
            />
            <BatchFile 
            :type="currentType"
            v-if="types.map(t => t.type).includes(currentType)"
            />
            <BatchProcess
            :process="processes.filter(p => p.key === currentType)[0]"
            v-if="processes.map(p => p.key).includes(currentType)"
            />
          </v-flex>
        </v-layout>
      </v-container>
    </div>
  </div>
</template>

<script>
import BatchFile from '@/components/BatchFile'
import BatchDate from '@/components/BatchDate'
import BatchProcess from '@/components/BatchProcess'

export default {
  props: ['batchKey'],
  data () {
    return {
      parameters: [
        {text: 'Date de début', key: 'dateDebut', prop: 'date_debut'},
        {text: 'Date de fin', key: 'dateFin', prop: 'date_fin'},
        {text: 'Date de fin effectifs', key: 'dateFinEffectif', prop: 'date_fin_effectif'}
      ],
      processes: [
        {text: 'Suppression',
          color: 'red',
          key: 'reset',
          img: '/static/poubelle.png',
          description: 'Retour au batch précédent',
          do (self) { self.$axios.get('/api/batch/reset') }
        },
        {text: 'Purger',
          color: 'blue',
          key: 'purge',
          img: '/static/gomme.svg',
          description: 'Retour au paramétrage',
          do (self) { self.$axios.get('/api/batch/purge') }
        },
        {text: 'Calcul Prédictions',
          color: 'green',
          key: 'predict',
          img: '/static/warning.png',
          description: 'Intégration des données et calcul des prédictions.',
          do (self) { self.$axios.get('/api/batch/process') }
        }
      ]
    }
  },
  methods: {
    setCurrentType (type) {
      this.currentType = type
    }
  },
  computed: {
    currentType: {
      get () { return this.$store.state.currentType },
      set (type) { this.$store.commit('setCurrentType', type) }
    },
    currentBatch () {
      return this.$store.state.batches.filter(b => b.id.key === this.batchKey)
    },
    types () {
      return this.$store.state.types.sort((a, b) => a.text.localeCompare(b.text))
    },
    features () {
      return this.$store.state.features
    },
    files () {
      return this.$store.state.files
    }
  },
  components: { BatchFile, BatchDate, BatchProcess }
}
</script>

<style>
  .selected {
    color: blue;
    font-size: 14px;
  }
  .container{
    display: flex;
  }
  .fixed{
    width: 300px;
  }
  .flex-item{
    flex-grow: 1;
  }
</style>