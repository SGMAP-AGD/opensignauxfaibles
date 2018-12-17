<template>
  <div >
    <v-card
    @click="hover=true"
      style="height: 80px; background: linear-gradient(#fff, #eee 45%, #ccc);"
      class="elevation-2 ma-2"
    >
      <div style="height: 100%; width: 100%;">
        <div class="entete">
          <PredictionWidgetScore id="widget" :prob="prediction.prob" :diff="prediction.diff"/>
        </div>
        <div class="corps">
          <div 
          style="left: 250px; position: absolute;"
          :id="'marge_' + prediction._id.siret"></div>
          <div style="white-space: nowrap; max-width: 400px; max-height:30px">
          <span style="font-size: 18px; color: #333; line-height: 10px; font-family: 'Oswald';">{{ prediction.etablissement.sirene.raisonsociale }}<br style="line-height: 10px;"/></span>
          </div>
          <span style="font-size: 12px; color: #333; line-height: 10px;">{{ prediction._id.siret }}<br style="line-height: 10px;"/></span>
          <v-img style="position: absolute; left: 160px; bottom: 10px;" width="17" src="/static/gray_apart.svg"></v-img>
          <v-img style="position: absolute; left: 90px; bottom: 10px;" width="57" src="/static/gray_urssaf.svg"></v-img>
          <div style="position: absolute; left: 195px; bottom: 4px; color: #333">
            <span style="font-size: 20px">{{ prediction.etablissement.effectif.effectif || 'n/c' }}</span><br/></div>
            <div v-if="hover"
            style="height: 50px; position: absolute; left: 400px; top: 5px; min-height: 100px; width: 100px;">
            Taux de marge
              <IEcharts
              style="height: 50px; background: #fff4"
                :loading="chart"
                :option="getMargeOption()"
              />
            </div>
            <div v-if="hover"
            style="height: 50px; position: absolute; left: 520px; top: 5px; min-height: 100px; width: 100px;">
            FRNG
              <IEcharts
              style="height: 50px; background: #fff4"
                :loading="chart"
                :option="getFRNGOption()"
              />
            </div>
            <div v-if="hover"
            style="height: 50px; position: absolute; left: 640px; top: 5px; min-height: 100px; width: 100px;">
            Financier CT
              <IEcharts
              style="height: 50px; background: #fff4"
                :loading="chart"
                :option="getFinCTOption()"
              />
            </div>
        </div>
        
      </div>
    </v-card>
  </div>
</template>

<script>
import 'echarts/lib/chart/line'
import 'echarts/lib/component/title'
import PredictionWidgetScore from '@/components/PredictionWidgetScore'

export default {
  props: ['prediction'],
  components: {
    PredictionWidgetScore
  },
  data () {
    return {
      hover: false,
      chart: false
    }
  },
  mounted () {
    // var paper = new Raphael('marge_' + this.prediction._id.siret, 20, 20)
    // var circle = paper.circle(10, 10, 10)
    // circle.attr('fill', '#00f')
    // circle.attr('stroke', '#fff')
  },
  methods: {

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
