<template>
  <div>
    <v-card
      style="height: 100px;"
      class="elevation-6 ma-3"
    >
      <v-container>
        <v-layout>
          <v-flex>
            {{ prediction._id.siret }}
          </v-flex>
          <v-flex>
            <PredictionWidgetScore :prob="prediction.prob" :diff="prediction.diff"/>
          </v-flex>
        </v-layout>
      </v-container>
    </v-card>
  </div>
</template>

<script>
import IEcharts from 'vue-echarts-v3/src/lite.js'
import 'echarts/lib/chart/line'
import 'echarts/lib/component/title'
import PredictionWidgetScore from '@/components/widgetPrediction'

export default {
  props: ['prediction'],
  components: {
    IEcharts,
    PredictionWidgetScore
  },
  beforeDestroy: function() {
    window.removeEventListener('resize', this.handleResize)
  },
  data: () => ({
  }),
  computed: {
    height() {
      return this.$store.state.height
    },
    scrollTop() {
      return this.$store.state.scrollTop
    },
    detectionLength() {
      return Math.round((this.height + this.scrollTop) / 1000 + 1)*10 ;
    }
  },
  methods: {
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
