<template>
  <div>
    <!-- <v-navigation-drawer
    class='elevation-6'
    absolute
    permanent
    :mini-variant = 'mini"
    style="z-index: 1"
    >

    </v-navigation-drawer>-->
    <!-- <v-data-table
    v-model="selected"
    :headers="headers"
    :items="predictionFiltered"
    :pagination.sync="pagination"
    select-all
    item-key="name"
    class="elevation-1"
    loading="false"
    :rows-per-page-items="[10]"
    >
      <template slot="headers" slot-scope="props">
        <tr>
          <th/>
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
        <tr
        :active="props.selected"
        >
          <td>
            <v-icon
            @click.left="open(props.item, true)"
            @click.middle="open(props.item, false)"
            >
              fa-address-card
            </v-icon>
          </td>
          <td ><v-tooltip left>
            <div slot="activator">{{ props.item.raison_sociale }}</div>
            {{ props.item._id.siret }}
            </v-tooltip> </td>
          <td class="text-xs-center"><widgetPrediction :prob="props.item.prob" :diff="props.item.diff"/></td>
          <td class="text-xs-right">
            {{ props.item.effectif }}
          </td>
          <td class="text-xs-right">{{ props.item.default_urssaf?"oui":"non" }}</td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.taux_marge}
                })
              )"
            />
          </td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.poids_frng}
                })
              )"
            />
          </td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.financier_court_terme}
                })
              )"
            />
          </td>
        </tr>
      </template>
    </v-data-table>-->
    <v-card
      style="height: 100px;"
      v-for="p in predictionFiltered"
      :key="p._id.siret"
      class="elevation-6"
    >
      <v-container>
        <v-layout>
          <v-flex>
          {{ p._id.siret }}
          </v-flex>
          <v-flex>

          <IEcharts
            style="height: 40px; width: 100px; "
            :loading="chart"
            :option="getMargeOption(
              (p.bdf || []).map(b => {
              return {'x': b.annee, 'y': b.taux_marge}
              })
            )"
          />

          </v-flex>
          <v-flex>

          <IEcharts
            style="height: 40px; width: 100px;"
            :loading="chart"
            :option="getMargeOption(
              (p.bdf || []).map(b => {
              return {'x': b.annee, 'y': b.poids_frng}
              })
            )"
          />

          </v-flex>
          <v-flex>

          <IEcharts
            style="height: 40px; width: 100px;"
            :loading="chart"
            :option="getMargeOption(
              (p.bdf || []).map(b => {
              return {'x': b.annee, 'y': b.financier_court_terme}
              })
            )"
          />

          </v-flex>
          <v-flex>

          <widgetPrediction :prob="p.prob" :diff="p.diff"/>

          </v-flex>

        </v-layout>
      </v-container>

    </v-card>
  </div>
</template>

<script>
import IEcharts from 'vue-echarts-v3/src/lite.js';
import 'echarts/lib/chart/line';
import 'echarts/lib/component/title';
import widgetPrediction from '@/components/widgetPrediction';
export default {
  props: ['batchKey'],
  components: {
    IEcharts,
    widgetPrediction
  },
  beforeDestroy: function() {
    window.removeEventListener('resize', this.handleResize);
  },
  data: () => ({
    effectifClass: [10, 20, 50, 100],
    selected: [],
    mini: true,
    loading: false,
    chart: false,
    pagination: {
      sortBy: 'name'
    },
    headers: [
      {
        text: 'raison sociale',
        align: 'left',
        value: 'raison_sociale'
      },
      { text: 'détection', value: 'prob' },
      { text: 'emploi', value: 'effectif' },
      { text: 'Défault urssaf', value: 'default_urssaf' },
      { text: 'Taux de marge', value: 'taux_marge' },
      { text: 'Fond de roulement', value: 'fond_roulement' },
      { text: 'Financier court terme', value: 'financier_court_terme' }
    ],
    prediction: [],
    naf: 'Industrie manufacturière',
    minEffectif: 20,
    entrepriseConnue: true,
    horsCCSF: true,
    horsProcol: true
  }),
  computed: {
    nomini: {
      get() {
        return !this.mini;
      },
      set(mini) {
        this.mini = !mini;
      }
    },
    predictionFiltered() {
      return this.prediction.slice(0, this.detectionLength);
    },
    tabs: {
      get() {
        return this.$store.getters.getTabs;
      },
      set(tabs) {
        this.$store.dispatch('updateTabs', tabs);
      }
    },
    activeTab: {
      get() {
        return this.$store.getters.activeTab;
      },
      set(activeTab) {
        this.$store.dispatch('updateActiveTab', activeTab);
      }
    },
    height() {
      return this.$store.state.height;
    },
    scrollTop() {
      return this.$store.state.scrollTop;
    },
    detectionLength() {
      return Math.round((this.height + this.scrollTop) / 1000 + 1)*10 ;
    }
  },
  mounted() {
    this.getPrediction();
  },
  methods: {
    open(etab, focus) {
      if (this.tabs.findIndex(t => t.siret === etab._id.siret) === -1) {
        let i = this.tabs.push({
          type: 'Etablissement',
          param: etab.raison_sociale,
          siret: etab._id.siret,
          batch: '1802'
        });
        if (focus) {
          this.activeTab = i - 1;
        }
      }
    },
    applyFilter(p) {
      return (
        (this.naf === 'Tous' || p.naf1 === this.naf) &&
        p.effectif >= this.minEffectif &&
        (p.connu === false || this.entrepriseConnue === false) &&
        (p.ccsf === false || this.horsCCSF === false) &&
        (p.procol === 'in_bonis' || this.horsProcol === false)
      );
    },
    getNAF() {
      var self = this;
      this.$axios.get('/api/data/naf').then(response => {
        self.naf = response.data;
      });
    },
    toggleAll() {
      if (this.selected.length) this.selected = [];
      else this.selected = this.desserts.slice();
    },
    changeSort(column) {
      if (this.pagination.sortBy === column) {
        this.pagination.descending = !this.pagination.descending;
      } else {
        this.pagination.sortBy = column;
        this.pagination.descending = false;
      }
    },
    getMargeOption(marge) {
      return {
        title: {
          text: null
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#283b56'
            }
          }
        },
        toolbox: {
          show: true
        },
        xAxis: {
          show: true,
          type: 'category',
          axisTick: false,
          data: marge.map(m => m.x)
        },
        yAxis: {
          type: 'value',
          show: false,
          min: -150,
          max: 150
        },
        series: [
          {
            color: 'indigo',
            smooth: true,
            name: 'taux marge',
            type: 'line',
            data: marge.map(m => m.y)
          }
        ]
      };
    },
    getPrediction() {
      this.loading = true;
      var self = this;
      this.$axios.post('/api/data/prediction').then(response => {
        var prediction = response.data;
        prediction.forEach(p => {
          p.bdf = Object.keys(p.bdf || {})
            .map(b => p.bdf[b])
            .sort((a, b) => a.annee < b.annee);
        });
        this.prediction = prediction;
        self.loading = false;
      });
    }
  }
};
</script>

<style scoped>
.echarts {
  width: 400px;
}

.pointer:hover {
  cursor: hand;
}
</style>
