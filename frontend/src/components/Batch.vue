<template>
<div>
  <v-card-text>
    <v-menu
      ref="menu_date_debut"
      :close-on-content-click="false"
      v-model="date_debut_status"
      :nudge-right="40"
      :return-value.sync="date_debut"
      lazy
      transition="scale-transition"
      offset-y
      max-width="290px"
      min-width="290px"
    >
      <v-text-field
        slot="activator"
        v-model="date_debut"
        label="Date de début"
        prepend-icon="event"
        readonly
      ></v-text-field>
      <v-date-picker
        v-model="date_debut"
        type="month"
        no-title
        scrollable
      >
        <v-spacer></v-spacer>
        <v-btn flat color="primary" @click="date_debut_status = false">Cancel</v-btn>
        <v-btn flat color="primary" @click="$refs.menu_date_debut.save(date_debut)">OK</v-btn>
      </v-date-picker>
    </v-menu>

    <v-menu
      ref="menu_date_fin"
      :close-on-content-click="false"
      v-model="date_fin_status"
      :nudge-right="40"
      :return-value.sync="date_fin"
      lazy
      transition="scale-transition"
      offset-y
      max-width="290px"
      min-width="290px"
      right
    >
      <v-text-field
        slot="activator"
        v-model="date_fin"
        label="Date de fin"
        prepend-icon="event"
        readonly
      ></v-text-field>
      <v-date-picker
        v-model="date_fin"
        type="month"
        no-title
        scrollable
      >
        <v-spacer></v-spacer>
        <v-btn flat color="primary" @click="date_fin_status = false">Cancel</v-btn>
        <v-btn flat color="primary" @click="$refs.menu_date_fin.save(date_fin)">OK</v-btn>
      </v-date-picker>
    </v-menu>

    <v-menu
      ref="menu_date_fin"
      :close-on-content-click="false"
      v-model="date_fin_effectif_status"
      :nudge-right="40"
      :return-value.sync="date_fin_effectif"
      lazy
      transition="scale-transition"
      offset-y
      max-width="290px"
      min-width="290px"
    >
      <v-text-field
        slot="activator"
        v-model="date_fin_effectif"
        label="Date de fin effectif"
        prepend-icon="event"
        readonly
      ></v-text-field>
      <v-date-picker
        v-model="date_fin_effectif"
        type="month"
        no-title
        scrollable
      >
        <v-spacer></v-spacer>
        <v-btn flat color="primary" @click="date_fin_effectif_status = false">Cancel</v-btn>
        <v-btn flat color="primary" @click="$refs.menu_date_fin.save(date_fin_effectif)">OK</v-btn>
      </v-date-picker>
    </v-menu>
  </v-card-text>
  <v-expansion-panel popout>
    <v-expansion-panel-content
      v-for="(files, type) in batch.files"
      v-bind:key="type"
    >
      <div slot="header">{{ type }} - {{ files.length }} {{ (files.length > 1) ? "fichiers":"fichier" }} </div>
      <v-card>
        <v-card-text v-for="file in files" v-bind:key="file">{{ file }}</v-card-text>
      </v-card>
    </v-expansion-panel-content>
  </v-expansion-panel>
</div>
</template>

<script>
export default {
  name: 'Batch',
  props: {
    batch: {
      type: Object,
      default: () => ({})
    }
  },
  data () {
    return {
      date_debut: this.batch.date_debut.substring(0, 7),
      date_debut_status: false,
      date_fin: this.batch.date_fin.substring(0, 7),
      date_fin_status: false,
      date_fin_effectif: this.batch.date_fin_effectif.substring(0, 7),
      date_fin_effectif_status: false
    }
  }
}
</script>
