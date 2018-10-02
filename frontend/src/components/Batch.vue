<template>
    <div >
      
    <v-navigation-drawer
    class="elevation-6"
    absolute
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
          :key="batchKey + param.key"
          ripple
          @click="setCurrentType(param.key)">
            <v-list-tile-content> 
              <v-list-tile-title
              :class="(param.key===currentType) ? 'selected': null"
              >{{ param.text }}</v-list-tile-title>
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
          <v-list-tile
          v-for="type in types"
          :key="type.text"
          @click="setCurrentType(type.type)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(type.type==currentType) ? 'selected': null">
              {{ type.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-divider></v-divider>
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
          :key="batchKey + process.key"
          @click="setCurrentType(process.key)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(process.key==currentType) ? 'selected': null">
              {{ process.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
      </v-list>
      
      <v-img src="/static/logo.png"> </v-img>
    </v-navigation-drawer>

    <div class="widget">
      <BatchDate 
      class="elevation-6"
      :key="batchKey + 'batchDate'"
      :date="currentType"
      :param="parameters.filter(p => p.key === currentType)[0]"
      v-if="parameters.map(p => p.key).includes(currentType)"
      />
      <BatchFile 
      :key="batchKey + 'batchFile'"
      :type="currentType"
      v-if="types.map(t => t.type).includes(currentType)"
      />
      <BatchProcess
      :key="batchKey + 'batchProcess'"
      :process="processes.filter(p => p.key === currentType)[0]"
      v-if="processes.map(p => p.key).includes(currentType)"
      />
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
      currentType: null,
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
  computed: {
    currentBatch () {
      return this.$store.state.batches.filter(b => b.id.key === this.batchKey)
    },
    features () {
      return this.$store.state.features
    },
    types () {
      return this.$store.state.types
    }
  },
  methods: {
    setCurrentType (type) {
      this.currentType = type
    }
  },
  components: { BatchFile, BatchDate, BatchProcess }
}
</script>

<style>
.selected {
  color: #700;
  font-size: 15px;
}
.widget {
  position: absolute;
  left: 320px;
  top: 20px; 
  right: 20px;
}
</style>