<script>
import axios from 'axios'
import { Line } from 'vue-chartjs'
export default {
  extends: Line,
  name: 'chartDebit',
  props: ['siret'],
  data () {
    return {
      datacollection: {
        datasets: []
      },
      options: {
        scales: {
          yAxes: [{
            stacked: false,
            ticks: {
              beginAtZero: true
            },
            gridLines: {
              display: true
            }
          }],
          xAxes: [{
            type: 'time',
            time: {
              unit: 'month'
            },
            gridLines: {
              display: false
            }
          }]
        },
        legend: {
          display: false
        },
        responsive: true,
        maintainAspectRatio: false
      }
    }
  },
  watch: {
    siret: function (newVal, oldVal) { // watch it
      if (newVal.length === 14 && !isNaN(newVal)) {
        this.draw()
      }
    }
  },
  mounted: {
    function () {
      if (this.siret.length === 14 && !isNaN(this.siret)) {
        this.draw()
      }
    }
  },
  methods: {
    draw () {
      var self = this
      self.datacollection.datasets = []
      axios.get(self.$api + '/data/debit/' + this.siret).then(function (response) {
        console.log(response.data)
        var debit = (response.data.value.batch['1802'].debit || {})
        var cotisation = (response.data.value.batch['1802'].cotisation || {})
        console.log(debit)
        // Mise en forme de la cotisation
        var calCot = Object.keys(cotisation).reduce(function (m, k) {
          var c = cotisation[k]
          var start = new Date(c.periode.start)
          var end = new Date(c.periode.end)
          var periods = self.$generatePeriodSerie(start, end)
          periods.map(function (date) {
            m[date] = (m[date] || 0) + c.du / periods.length
          })
          return m
        }, {})
        var keysCot = Object.keys(calCot)
        keysCot.sort(function (a, b) {
          a = new Date(a)
          b = new Date(b)
          return a > b ? 1 : a < b ? -1 : 0
        })

        var datasetCot = keysCot.map(function (date) {
          return {t: date, y: calCot[date]}
        })
        self.datacollection.datasets.push({
          label: 'Cotisation',
          radius: 0,
          borderWidth: 1,
          steppedLine: true,
          backgroundColor: 'rgba(0,0,128,0.4)',
          data: datasetCot})

        // Mise en forme de la dette cumulée
        var calDette = {}

        Object.keys(debit).map(function (h) {
          var d = debit[h]
          var start = d.periode.start
          var end = d.periode.end
          var numEcn = d.numero_ecart_negatif
          var compte = d.numero_compte
          var key = start + '-' + end + '-' + numEcn + '-' + compte
          var dateTraitement = new Date(d.date_traitement)
          calDette[dateTraitement] = (calDette[dateTraitement] || {})
          calDette[dateTraitement][key] = (calDette[dateTraitement][key] || {'numero_historique': 0})
          calDette[dateTraitement][key] = (calDette[dateTraitement][key].numero_historique > d.numero_historique ? calDette[dateTraitement][key] : d)
        })

        var keys = Object.keys(calDette)
        keys.sort(function (a, b) {
          a = new Date(a)
          b = new Date(b)
          return a > b ? 1 : a < b ? -1 : 0
        })
        var detteCumulee = []
        keys.reduce(function (m, k) {
          Object.keys(calDette[k]).map(function (key) {
            m[key] = calDette[k][key]
          })
          detteCumulee[k] = Object.keys(m).reduce(function (total, partie) {
            return total + m[partie].part_patronale + m[partie].part_ouvriere
          }, 0)
          return m
        }, {})

        var dataset = keys.map(function (k) {
          return {t: k, y: Math.round(detteCumulee[k])}
        })
        self.datacollection.datasets.push({
          label: 'Dette Cumulée',
          radius: 2,
          borderWidth: 3,
          steppedLine: true,
          backgroundColor: 'rgba(128,0,0,0.4)',
          data: dataset})
        self.renderChart(self.datacollection, self.options)
      })
    }
  }
}
</script>
