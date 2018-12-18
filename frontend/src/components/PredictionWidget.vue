<template>
  <div >
    <v-card
      @click="showEtablissement()"
      style="height: 80px; background: linear-gradient(#fff, #eee 45%, #ccc);"
      class="elevation-2 ma-2 pointer"
    >
      <div style="height: 100%; width: 100%;" >
        <div class="entete pointer" >
          <PredictionWidgetScore id="widget" :prob="prediction.prob" :diff="prediction.diff"/>
        </div>
        <div class="corps">
          <div 
          style="left: 250px; position: absolute;"
          :id="'marge_' + prediction._id.siret"></div>
          <div style="white-space: nowrap; overflow: hidden; max-width: 400px; max-height:30px">
            <span style="font-size: 18px; color: #333; line-height: 10px; font-family: 'Oswald';">{{ prediction.etablissement.sirene.raisonsociale }}<br style="line-height: 10px;"/></span>
          </div>
          <div style="left: 500px; position: absolute; top: 4px">
            Chiffre d'affaire (k€)<br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/a' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].ca || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].ca || 'n/a' }}</td>
              </tr>
            </table>
          </div>

          <div style="left: 500px; position: absolute; top: 4px">
            EBE (k€)<br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/a' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].excedent_brut_d_exploitation || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].excedent_brut_d_exploitation || 'n/a' }}</td>
              </tr>
            </table>
          </div>

          <div style="left: 500px; position: absolute; top: 4px">
            Bénéfice ou perte (k€)<br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/a' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].benefice_ou_perte || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].benefice_ou_perte || 'n/a' }}</td>
              </tr>
            </table>
          </div>

          <div style="left: 500px; position: absolute; top: 4px">
            Bénéfice ou perte (k€)<br/>
            <table>
              <tr>
                <td>{{ prediction.entreprise.diane[0].exercice_diane || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].exercice_diane || 'n/a' }}</td>
              </tr>
              <tr>
                <td>{{ prediction.entreprise.diane[0].liquidite_reduite || 'n/a' }}</td>
                <td>{{ prediction.entreprise.diane[1].liquidite_reduite || 'n/a' }}</td>
              </tr>
            </table>
          </div>

          <span style="font-size: 12px; color: #333; line-height: 10px;">{{ prediction._id.siret }}<br style="line-height: 10px;"/></span>
          <v-img style="position: absolute; left: 160px; bottom: 10px;" width="17" src="/static/gray_apart.svg"></v-img>
          <v-img style="position: absolute; left: 90px; bottom: 10px;" width="57" :src="'/static/' + (prediction.etablissement.urssaf?'red':'gray') + '_urssaf.svg'"></v-img>
          <div style="position: absolute; left: 195px; bottom: 4px; color: #333">
            <span :class="variationEffectif" style="font-size: 20px">{{ prediction.etablissement.effectif || 'n/c' }}</span>
          </div>

          {{ prediction.entreprise.diane[0].ca }} {{ prediction.entreprise.diane[1].ca }} <br/>
          {{ prediction.entreprise.diane[0].benefice_ou_perte }} {{ prediction.entreprise.diane[1].benefice_ou_perte }} <br/>
          {{ prediction.entreprise.diane[0].liquidite_reduite }} {{ prediction.entreprise.diane[1].liquidite_reduite }} <br/>


        </div>
        <v-dialog 
        attach="#detection"
        lazy
        fullscreen
        v-model="dialog">
          <div style="height: 100%; width: 100%;  font-weight: 800; font-family: 'Abel', sans;">
            <v-toolbar fixed class="toolbar" height="35px" style="color: #fff; font-size: 22px;">
              <v-spacer/>
                {{ prediction.etablissement.sirene.raisonsociale }}
              <v-spacer/>
              <v-icon @click="dialog=false"  style="color: #fff">mdi-close</v-icon>
            </v-toolbar>
          <Etablissement :siret="prediction._id.siret"></Etablissement>
          </div>
        </v-dialog>
      </div>
    </v-card>

  </div>
</template>

<script>
import Etablissement from '@/components/Etablissement'
import PredictionWidgetScore from '@/components/PredictionWidgetScore'

export default {
  props: ['prediction'],
  components: {
    PredictionWidgetScore,
    Etablissement
  },
  data () {
    return {
      dialog: false
    }
  },
  computed: {
    variationEffectif () {
      if (this.prediction.etablissement.effectif / this.prediction.etablissement.effectif_precedent > 1.05) {
        return 'high'
      }
      if (this.prediction.etablissement.effectif / this.prediction.etablissement.effectif_precedent < 0.95) {
        return 'down'
      }
      return 'none'
    }
  },
  methods: {
    showEtablissement () {
      this.dialog = true
      console.log('yeah')
    }
  }
}
</script>

<style scoped>
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
.high {
  color: rgb(16, 114, 16);
}
.down {
  color: rgb(139, 19, 19);
}
td {
  width: 60px;
}
.pointer {cursor: pointer;}
</style>
