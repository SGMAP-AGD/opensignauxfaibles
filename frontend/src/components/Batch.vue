<template>
<div>
  <v-card-text>
    Batch {{batch.id.key}}
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
  }
}
</script>
