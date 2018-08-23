<template>
  <v-container fluid grid-list-md text-xs-center>
      <v-layout row wrap justify-start>
        <v-flex xs12>TABLEAU DE BORD</v-flex>
        
        <v-flex xs4>
          <v-card>
            <v-card-title>
              <h2>Intégration des données</h2>
            </v-card-title>
                        <v-divider></v-divider>
            <v-card-text>
              <v-list>
                <v-list-tile avatar>
                  Chiffre 1
                </v-list-tile>
                <v-list-tile avatar>
                  Chiffre 2
                </v-list-tile>
              </v-list>
            </v-card-text>
          </v-card>
        </v-flex>
        <v-flex xs4>
          <v-card>
            <v-card-title>
              <h2>Calcul prédictions</h2>
            </v-card-title>
            <v-card-text>
              <v-list>
                <v-list-tile avatar>
                  Effectué: x% 
                </v-list-tile>

                <v-list-tile avatar>
                    En attente: y%
                </v-list-tile>
              </v-list>
            </v-card-text>
          </v-card>
        </v-flex>
        <v-flex xs4>
          <v-card>
            <v-card-title>
              <h2>Visites d'entreprises</h2>
            </v-card-title>
            <v-card-text>
              <v-list>
                <v-list-tile avatar>
                  Planifiées: x% 
                </v-list-tile>  
                <v-list-tile avatar>
                  Effectuées: y%
                </v-list-tile>
                <v-list-tile avatar>
                  Retard: y%
                </v-list-tile>
              </v-list>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    
  </v-container>
</template>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
