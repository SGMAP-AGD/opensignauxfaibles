package main

import (
	"fmt"

	"github.com/chrnin/ganboard"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Lancer Rserve en background

	loadConfig()

	go r()

	r := gin.Default()
	r.Use(DB())
	r.Use(Kanboard())
	// FIXME: configurer correctement CORS
	r.Use(cors.Default())

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	api := r.Group("api")
	{
		api.OPTIONS("/auth", auth)
		api.GET("/admin/region", AdminRegion)
		api.POST("/admin/region", AdminRegionAdd)
		api.DELETE("/admin/region", AdminRegionDelete)
		api.GET("/repo/create/:region/:periode", createRepo)
		api.GET("/purge", purge)

		api.GET("/kanboard/task/create/:siret", getKBProject)

		api.GET("/import/bdf/:region/:batch", importBDF)

		api.GET("/import/all/:region/:batch", importAll)
		api.GET("/import/apdemande/:region/:batch", importAPDemande)
		api.GET("/import/apconso/:region/:batch", importAPConso)
		api.GET("/import/cotisation/:region/:batch", importCotisation)
		api.GET("/import/debit/:region/:batch", importDebit)
		api.GET("/import/effectif/:region/:batch", importEffectif)
		api.GET("/import/altares/:region/:batch", importAltares)
		api.GET("/import/delai/:region/:batch", importDelai)
		api.GET("/import/sirene/:region/:batch", importSirene)
		api.GET("/compact/etablissement/:siret", compact)
		api.GET("/compact/etablissement", compactAll)
		api.GET("/compact/entreprise/:siren", compactEntreprise)
		api.GET("/compact/entreprise", compactAllEntreprise)

		api.GET("/reduce/:siret", reduce)

		api.GET("/reduce", reduceAll)
		api.GET("/etablissement/:siret", browseEtablissement)
		api.GET("/orig/:siret", browseOrig)
		api.POST("/R/algo1", algo1)
		api.GET("/data/debit/:siret", dataDebit)
	}
	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/opensignauxfaibles")
	viper.AddConfigPath("$HOME/.opensignauxfaibles")
	viper.AddConfigPath(".")
	viper.SetDefault("APP_BIND", ":3000")
	viper.SetDefault("APP_DATA", "$HOME/data-raw/")
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "27017")
	viper.SetDefault("DB", "opensignauxfaibles")
	viper.SetDefault("JWT_SECRET", "One might change this because one day it will not be sufficient")
	viper.SetDefault("KANBOARD_ENDPOINT", "http://localhost/kanboard/jsonrpc.php")
	viper.SetDefault("KANBOARD_USERNAME", "admin")
	viper.SetDefault("KANBOARD_PASSWORD", "admin")
	err := viper.ReadInConfig()
	fmt.Println(err)
}

// Kanboard -> ganboard.Client in gingonic context
func Kanboard() gin.HandlerFunc {

	endpoint := viper.GetString("KANBOARD_ENDPOINT")
	username := viper.GetString("KANBOARD_USERNAME")
	password := viper.GetString("KANBOARD_PASSWORD")

	client := ganboard.Client{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}

	return func(c *gin.Context) {
		c.Set("KBCLIENT", &client)
		c.Next()
	}
}

func debugalgo1(c *gin.Context) {

}
