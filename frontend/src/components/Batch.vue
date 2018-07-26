<template>
    <v-card
      raised>
      
        <v-toolbar 
          card
          :color="(batch.readonly)?'deep-purple lighten-4':(batch.altered)?'red lighten-4':(newBatch)?'orange lighten-4':'teal lighten-4'">
          <v-menu 
            offset-y
            v-if="!newBatch">
            <v-btn
              :color="(batch.readonly)?'deep-purple lighten-2':(batch.altered)?'red darken-1':(newBatch)?'orange darken-1':'teal darken-1'"
              dark
              slot="activator">
              {{ newBatch ? "" : batch.id.key }}</v-btn>
            <v-list>
              <v-list-tile 
                key="drop"
                @click="dropBatch(batch)"
                >
                <v-list-tile-action>
                  <v-icon>fa-trash</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>Suppression</v-list-tile-content>
              </v-list-tile>
              <v-list-tile 
                key="delete"
                @click="test()"
                >
                <v-list-tile-action>
                  <v-icon>fa-eraser</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>Purge</v-list-tile-content>
              </v-list-tile>
              <v-list-tile 
                key="purge"
                @click="test()"
                >
                <v-list-tile-action>
                  <v-icon>fa-download</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>Import</v-list-tile-content>
              </v-list-tile>

              <v-list-tile 
                @click="test()"
                key="reduce"
                >
                <v-list-tile-action>
                  <v-icon>fa-calculator</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>Réduction</v-list-tile-content>
              </v-list-tile>
              <v-list-tile 
                @click="test()"
                key="predict"
                >
                <v-list-tile-action>
                  <v-icon>fa-registered</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>Prédiction</v-list-tile-content>
              </v-list-tile>
            </v-list>
          </v-menu>
          <v-container>
            <v-layout>
              <v-flex>
              <v-text-field
                
                dense
                label="Nouveau lot"
                v-if="newBatch"
                v-bind:key="'newBatch'"
                v-model="batch.id.key"
                color="pink darken-4"
                @input="alterBatch()"
                :rules="[duplicateBatchKey,newBatchKey]"
              ></v-text-field>
              </v-flex>
                          </v-layout>
          </v-container>
              <v-spacer></v-spacer>
              <v-flex xs2>
              <v-btn 
               v-if="!newBatch"
                icon
                :color="(batch.readonly)?'deep-purple lighten-2':(batch.altered)?'red darken-1':'teal darken-1'"
                :disabled="(!batch.altered || batch.readonly) && !newBatch || batch.id.key=='' || duplicateBatchKey() != true || newBatchKey() != true"
                :flat="(!batch.altered || batch.readonly) && !newBatch"
                raised
                :dark="batch.altered||newBatch"
                @click="saveBatch(batch)"
              >
                <v-icon light>
                  {{ (batch.readonly)?'fa-lock':(newBatch)?'fa-plus':(batch.altered)?'fa-save':'fa-lock-open' }}
                </v-icon>
              </v-btn>
              </v-flex>

        </v-toolbar>
        <v-divider></v-divider>
          <v-btn
          v-if="newBatch"
          depressed
          fab
          large
          @click="saveBatch(batch)">
            <v-icon>fa-plus</v-icon>
          </v-btn>
          <v-list dense v-if="!newBatch">
            <v-list-tile>
              <v-list-tile-content>Fichiers</v-list-tile-content>
                <v-list-tile-content class="align-end">
                <v-dialog
                  transition="dialog-bottom-transition"
                  lazy
                  scrollable
                  v-model="batch.dialog"
                  :disabled="batch.readonly"
                  width="80%"
                >
                  <v-btn 
                    flat
                    slot="activator"
                  >
                    {{ Object.keys(batch.files).reduce((a,b) => { return a += batch.files[b].length }, 0) }}
                  </v-btn>
                  <BatchFile :batch="batch" :types="types" :files="files"></BatchFile>
                </v-dialog>
              </v-list-tile-content>
            </v-list-tile>
            <v-list-tile>
                <v-list-tile-content>Date de Début</v-list-tile-content>
                <v-list-tile-content class="align-end">
                  <v-menu
                  :disabled="batch.readonly">
                    <v-btn
                      slot="activator"
                      readonly
                      flat
                      >{{ batch.params.date_debut.substring(0, 7) }}
                    </v-btn>
                  <v-date-picker 
                    header="false"
                    locale="fr-FR" 
                    type="month"
                    v-model="batch.params.date_debut"
                    @input="alterBatch(batch)"
                    :landscape="false" 
                    :reactive="false">
                  </v-date-picker>
                </v-menu>
                </v-list-tile-content>
            </v-list-tile>
            <v-list-tile>
                <v-list-tile-content>Date de Fin</v-list-tile-content>
                <v-list-tile-content class="align-end">
                  <v-menu
                    :disabled="batch.readonly"
                  >
                    <v-btn
                      slot="activator"
                      readonly
                      flat
                    >{{ batch.params.date_fin.substring(0, 7) }}</v-btn>
                  <v-date-picker 
                    header="false"
                    locale="fr-FR" 
                    type="month"
                    v-model="batch.params.date_fin"
                    @input="alterBatch(batch)"
                    :landscape="false" 
                    :reactive="false">
                  </v-date-picker>
                </v-menu>
                </v-list-tile-content>
            </v-list-tile>
            <v-list-tile>
                <v-list-tile-content>Derniers effectifs</v-list-tile-content>
                <v-list-tile-content class="align-end">
                  <v-menu
                  :disabled="batch.readonly">
                    <v-btn
                      slot="activator"
                      readonly
                      flat
                    >{{ batch.params.date_fin_effectif.substring(0, 7) }}</v-btn>
                  <v-date-picker 
                    header="false"
                    locale="fr-FR" 
                    type="month"
                    v-model="batch.params.date_fin_effectif"
                    @input="alterBatch(batch)"
                    :landscape="false" 
                    :reactive="false">
                  </v-date-picker>
                </v-menu>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile>
            
            </v-list-tile>
          </v-list>
    </v-card>
</template>

<script>
import axios from 'axios'
import BatchFile from '@/components/BatchFile'

export default {
  name: 'Batch',
  components: { BatchFile },
  methods: {
    test () {
      alert(null)
    },
    newBatchKey () {
      if (this.$parent.batches.filter(batch => batch.readonly).filter(batch => batch.id.key >= this.batch.id.key).length > 0 && this.batch.id.key !== '') {
        return this.$parent.batches.filter(batch => batch.readonly).filter(batch => batch.id.key >= this.batch.id.key).length + ' conflits.'
      } else {
        return true
      }
    },
    duplicateBatchKey () {
      if (this.$parent.batches.filter(batch => batch.id.key === this.batch.id.key).length) {
        return 'Doublon'
      } else {
        return true
      }
    },
    alterBatch () {
      this.batch.altered = true
    },
    lockBatch () {
      this.batch.readonly = true
    },
    unlockBatch () {
      this.batch.readonly = false
    },
    dropBatch (batch) {
      axios.delete(this.$api + '/admin/batch/' + batch.id.key, {}).then(this.$parent.refresh())
    },
    saveBatch (batch) {
      batch.params.date_debut = batch.params.date_debut + '-01T00:00:00Z'
      batch.params.date_fin = batch.params.date_fin + '-01T00:00:00Z'
      batch.params.date_fin_effectif = batch.params.date_fin_effectif + '-01T00:00:00Z'

      axios.post(this.$api + '/admin/batch', batch).then(response => {
        this.$parent.refresh()
        this.batch.altered = false
      })
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
    },
    newBatch: {
      type: Boolean,
      default: () => (false)
    }
  }
}
</script>
