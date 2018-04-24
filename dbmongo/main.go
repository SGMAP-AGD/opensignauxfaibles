package main

import (
	"fmt"

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
		api.GET("/import/all/:region/:batch", importAll)
		api.GET("/import/apdemande/:region/:batch", importAPDemande)
		api.GET("/import/apconso/:region/:batch", importAPConso)
		api.GET("/import/cotisation/:region/:batch", importCotisation)
		api.GET("/import/debit/:region/:batch", importDebit)
		api.GET("/import/effectif/:region/:batch", importEffectif)
		api.GET("/import/altares/:region/:batch", importAltares)
		api.GET("/import/delais/:region/:batch", importEffectif)

		api.GET("/import/compact/:siret", compact)
		api.GET("/import/compact", compactAll)
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
	err := viper.ReadInConfig()
	fmt.Println(err)

}
