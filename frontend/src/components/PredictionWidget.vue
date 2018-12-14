<template>
  <div>
    <v-card
      style="height: 80px; background: linear-gradient(#fff, #eee 45%, #ccc);"
      class="elevation-2 ma-2"
    >
    <div style="height: 100%; width: 100%;">
      <div class="entete">
        <PredictionWidgetScore id="widget" :prob="prediction.prob" :diff="prediction.diff"/>
      </div>
      <div class="corps">
        <span style="font-size: 10px">{{ prediction._id.siret }} effectif: {{ prediction.etablissement.effectif.effectif || 'n/c' }}</span><br/>
        <span style="font-size: 13px">{{ prediction.etablissement.sirene.raisonsociale }}</span><br/>
        <v-img style="position: absolute; right: 70px; bottom: 10px;" width="18" src="/static/red_apart.svg"></v-img>
        <v-img style="position: absolute; right: 10px; bottom: 10px;" width="50" src="/static/red_urssaf.svg"></v-img>
      </div>      
    </div>
    </v-card>
  </div>
</template>

<script>
// import IEcharts from 'vue-echarts-v3/src/lite.js'
import 'echarts/lib/chart/line'
import 'echarts/lib/component/title'
import PredictionWidgetScore from '@/components/widgetPrediction'

export default {
  props: ['prediction'],
  components: {
    // IEcharts,
    PredictionWidgetScore
  },
  methods: {
    getMargeOption (marge) {
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
      }
    }
  }
}
</script>

<style scoped>
.echarts {
  width: 400px;
}
div.entete {
  float: left;
  background: linear-gradient(270deg, rgba(119, 122, 170, 0.219), rgba(119, 122, 170, 0));
  border-right: solid 1px #3334;
  width: 80px;
  height: 80px;
  text-align: center;
  padding: 20px;
}
div.corps {
  flex: 1;
  padding: 5px;
  margin-left: 80px;
  height: 80px;
  background: linear-gradient(45deg, rgba(50, 51, 121, 0.212), #0000);
}
.pointer:hover {
  cursor: hand;
}
</style>
