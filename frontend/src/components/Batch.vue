<template>
  <div>
    <v-card
    raised>
      <v-card-title>
        <v-toolbar 
          flat
          :color="(batch.readonly)?'indigo lighten-5':(altered)?'red lighten-5':'light-green accent-1'">
        <v-menu 
          offset-y
          v-if="!newBatch">
          <v-btn 
            :color="(batch.readonly)?'indigo darken-4':(altered)?'red darken-4':'light-green darken-4'"
            dark
            slot="activator">
            {{ newBatch ? "" : batch.id.key }}</v-btn>
          <v-list>
            <v-list-tile 
              key="delete"
              @click="dropBatch(batch)"
              >
              <v-list-tile-action>
                <v-icon>fa-trash</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>Suppression</v-list-tile-content>
            </v-list-tile>
            <v-list-tile 
              key="purge"
              @click="test()"
              >
              <v-list-tile-action>
                <v-icon>fa-eraser</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>Purge</v-list-tile-content>
            </v-list-tile>
            <v-list-tile 
              key="import"
              @click="test()"
              >
              <v-list-tile-action>
                <v-icon>fa-download</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>Import</v-list-tile-content>
            </v-list-tile>
            <v-list-tile 
              key="import"
              @click="test()"
              
              >
              <v-list-tile-action>
                <v-icon>fa-compress</v-icon>
              </v-list-tile-action>
              <v-list-tile-content>Compactage</v-list-tile-content>
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
        
        <v-spacer ></v-spacer>
        <v-text-field
          v-if="newBatch"
          v-model="batch.id.key"
          @input="alterBatch()"
        ></v-text-field>
        <v-spacer></v-spacer>
          <v-btn 
            icon 
            :color="(batch.readonly)?'indigo darken-4':(altered)?'red darken-4':'light-green darken-4'"
            :disabled="(!altered || batch.readonly) && !newBatch || batch.id.key==''"
            :flat="(!altered || batch.readonly) && !newBatch"
            raised
            :dark="altered||newBatch"
            @click="saveBatch(batch)"
          >
            <v-icon light>
              {{ (batch.readonly)?'fa-lock':(newBatch)?'fa-plus':(altered)?'fa-save':'fa-lock-open' }}
            </v-icon>
          </v-btn>
          <v-tooltip top>
            coucou
            <span>Tu aimes les tooltips</span>
          </v-tooltip>    
          </v-toolbar>
      </v-card-title>
        <v-divider></v-divider>
          <v-list dense>
            <v-list-tile>
              <v-list-tile-content>Fichiers</v-list-tile-content>
                <v-list-tile-content class="align-end">
                <v-dialog
                :disabled="batch.readonly"
                  width="70%"
                >
                  <v-btn flat
                    slot="activator"
                  >
                    {{ Object.keys(batch.files).reduce((a,b) => { return a += batch.files[b].length }, 0) }}
                  </v-btn>
                  <BatchFile :batch="batch"></BatchFile>
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
                      
                    >{{ batch.params.date_debut.substring(0, 7) }}</v-btn>
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
  </div>
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
    alterBatch (batch) {
      this.altered = true
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
        this.altered = false
      })
    }
  },
  props: {
    batch: {
      type: Object,
      default: () => ({})
    },
    newBatch: {
      type: Boolean,
      default: () => (false)
    }
  },
  data () {
    return {
      altered: false

    }
  }
}
</script>
