<template>
  <div>
    <v-navigation-drawer permanent app right clipped mini-variant>
        gnagnagna
    </v-navigation-drawer>
    <v-data-table
        v-model="selected"
        :headers="headers"
        :items="prediction"
        :pagination.sync="pagination"
        select-all
        item-key="name"
        class="elevation-1"
    >
      <template slot="headers" slot-scope="props">
        <tr>
            <th>
            <v-checkbox
                :input-value="props.all"
                :indeterminate="props.indeterminate"
                primary
                hide-details
                @click.native="toggleAll"
            ></v-checkbox>
            </th>
            <th
            v-for="header in props.headers"
            :key="header.text"
            :class="['column sortable', pagination.descending ? 'desc' : 'asc', header.value === pagination.sortBy ? 'active' : '']"
            @click="changeSort(header.value)"
            >
            <v-icon small>arrow_upward</v-icon>
            {{ header.text }}
            </th>
        </tr>
      </template>
      <template slot="items" slot-scope="props">
        <tr :active="props.selected" @click="props.selected = !props.selected">
            <td>
            <v-checkbox
                :input-value="props.selected"
                primary
                hide-details
            ></v-checkbox>
            </td>
            <td>{{ props.item.siret }}</td>
            <td class="text-xs-right">{{ Math.round(props.item.score*1000)/1000 }}</td>
            <td class="text-xs-right">{{ props.item.fat }}</td>
            <td class="text-xs-right">{{ props.item.carbs }}</td>
            <td class="text-xs-right">{{ props.item.protein }}</td>
            <td class="text-xs-right">{{ props.item.iron }}</td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script>
  import axios from 'axios'

  export default {
    data: () => ({
      pagination: {
        sortBy: 'name'
      },
      selected: [],
      actual_batch: '1803',
      headers: [
        {
          text: 'siret',
          align: 'left',
          value: 'siret'
        },
        { text: 'score', value: 'score' },
        { text: 'Fat (g)', value: 'fat' },
        { text: 'Carbs (g)', value: 'carbs' },
        { text: 'Protein (g)', value: 'protein' },
        { text: 'Iron (%)', value: 'iron' }
      ],
      prediction: []
    }),
    mounted () {
      this.getPrediction()
    },
    methods: {
      toggleAll () {
        if (this.selected.length) this.selected = []
        else this.selected = this.desserts.slice()
      },
      changeSort (column) {
        if (this.pagination.sortBy === column) {
          this.pagination.descending = !this.pagination.descending
        } else {
          this.pagination.sortBy = column
          this.pagination.descending = false
        }
      },
      getPrediction () {
        var self = this
        axios.get(this.$api + '/data/prediction/1803/algo1/0').then(response => {
          self.prediction = response.data.map(prediction => {
            var a = self.projectBatch(prediction.etablissement)
            console.log(a)
            return {
              'siret': prediction._id.siret,
              'score': prediction.score
            }
          })
        })
      },
      projectBatch (o) {
        return Object.keys((o.batch || {})).sort()
          .filter(batch => batch <= this.actual_batch).reduce((m, batch) => {
            Object.keys(o.batch[batch]).forEach((type) => {
              m[type] = (m[type] || {})
              var arrayDelete = (o.batch[batch].compact.delete[type] || [])
              if (arrayDelete !== {}) {
                arrayDelete.forEach(hash => {
                  delete m[type][hash]
                })
              }
              Object.assign(m[type], o.batch[batch][type])
            })
            return m
          }, {})
      },
      flattenTypes (o) {
        return Object.keys(o).reduce((accu, type) => {
          accu[type] = Object.values(o[type])
        }, {})
      }
    }
  }
</script>