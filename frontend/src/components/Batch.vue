<template>
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
        { text: 'Date de début', key: 'dateDebut', prop: 'date_debut' },
        { text: 'Date de fin', key: 'dateFin', prop: 'date_fin' },
        { text: 'Date de fin effectifs', key: 'dateFinEffectif', prop: 'date_fin_effectif' }
      ],
      processes: [
        { text: 'Suppression',
          color: 'red',
          key: 'reset',
          img: '/static/poubelle.png',
          description: 'Retour au batch précédent',
          do (self) { self.$axios.get('/api/batch/revert') }
        },
        { text: 'Purger',
          color: 'blue',
          key: 'purge',
          img: '/static/gomme.svg',
          description: 'Retour au paramétrage',
          do (self) { self.$axios.get('/api/batch/purge') }
        },
        { text: 'Calcul Prédictions',
          color: 'green',
          key: 'predict',
          img: '/static/warning.png',
          description: 'Intégration des données et calcul des prédictions.',
          do (self) { self.$axios.get('/api/batch/process') }
        },
        { text: 'Clôture',
          color: 'black',
          key: 'close',
          img: '/static/warning.png',
          description: 'Clôture du batch en cours et création du suivant',
          do (self) { self.$axios.get('/api/batch/next') }
        }
      ]
    }
  },
  computed: {
    currentType: {
      get () {
        return this.$store.state.currentType
      },
      set (type) {
        this.$store.dispatch('setCurrentType', type)
      }
    },
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch: {
      get () {
        if (this.$store.state.batches !== []) {
          return this.$store.getters.batches[this.currentBatchKey]
        } else {
          return null
        }
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch).then(r => this.$store.dispatch('checkEpoch'))
      }
    },
    rightDrawer: {
      get () {
        return this.$store.state.rightDrawer
      },
      set (val) {
        this.$store.dispatch('setRightDrawer', val)
      }
    },
    features () {
      return this.$store.state.features
    },
    types () {
      return this.$store.state.types
    }
  },
  methods: {
    toggleComplete (type) {
      let batch = this.currentBatch
      if (batch.complete_types.includes(type)) {
        batch.complete_types = batch.complete_types.filter(t => t !== type)
      } else {
        batch.complete_types = (batch.complete_types || []).concat(type)
      }
      this.currentBatch = batch
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
</style>
