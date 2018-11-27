import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'
import 'vuetify/src/stylus/app.styl'
import fr from 'vuetify/es5/locale/fr'
import colors from 'vuetify/es5/util/colors'

Vue.use(Vuetify, {
  iconfont: 'md',
  lang: {
    locales: { fr },
    current: 'fr'
  },
  theme: {
    primary: '#20459a',
    secondary: '#8e0000',
    accent: colors.red.base
  }
})
