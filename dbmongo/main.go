package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Lancer Rserve en background

	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

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

		api.GET("/purge", purge)

		api.GET("/kanboard/get/projects", listProjects)
		api.GET("/kanboard/get/tasks", getKBTasks)

		api.GET("/import/all/:batch", importAll)
		api.GET("/import/apdemande/:batch", importAPDemande)
		api.GET("/import/apconso/:batch", importAPConso)
		api.GET("/import/cotisation/:batch", importCotisation)
		api.GET("/import/ccsf/:batch", importCCSF)
		api.GET("/import/debit/:batch", importDebit)
		api.GET("/import/effectif/:batch", importEffectif)
		api.GET("/import/altares/:batch", importAltares)
		api.GET("/import/delai/:batch", importDelai)
		api.GET("/import/sirene/:batch", importSirene)
		api.GET("/import/bdf/:batch", importBDF)
		api.GET("/import/dpae/:batch", importDPAE)

		api.GET("/repo/create/:batch", createRepo)

		api.GET("/compact/:siret", compactEtablissement)
		api.GET("/compact", compactEtablissement)

		api.GET("/reduce/:siren", reduce)

		api.GET("/reduce", reduceAll)
		api.GET("/browse/:siren", browse)
		api.GET("/orig/:siret", browseOrig)
		api.POST("/R/algo1", algo1)
		api.GET("/data/debit/:siret", dataDebit)

		api.GET("/debug/:param", debug)
	}
	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func debug(c *gin.Context) {
	param := c.Params.ByName("param")
	spew.Dump(batchToTime(param))
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
