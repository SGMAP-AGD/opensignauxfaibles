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

		api.PUT("/admin/batch/:batchID", registerNewBatch)

		api.GET("/admin/batch", listBatch)
		api.GET("/admin/files", adminFiles)
		api.POST("/admin/attach", attachFileBatch)
		api.GET("/admin/types", listTypes)
		api.GET("/admin/clone/:to", cloneDB)

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

		api.GET("/compact/etablissement/:siret", compactEtablissement)
		api.GET("/compact/etablissement", compactEtablissement)
		api.GET("/compact/entreprise/:siren", compactEntreprise)
		api.GET("/compact/entreprise", compactEntreprise)
		api.GET("/index/entreprise/:siren", indexEntreprise)

		api.GET("/reduce/:algo/:batch/:siret", reduce)
		api.GET("/reduce/:algo/:batch", reduce)

		api.POST("/data/browse", browse)
		api.POST("/R/algo1", algo1)

		api.POST("/data", data)
		api.GET("/data/batch", dataBatch)
		api.GET("/data/algo", dataAlgo)
		api.GET("/data/prediction/:batch/:algo/:page", predictionBrowse)

		api.GET("/debug/:batch", debug)
	}
	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func debug(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	batch := c.Params.ByName("batch")
	files, err := GetFileList(basePath, batch)
	spew.Dump(files)
	spew.Dump(err)
	c.JSON(200, "cooool")
	c.JSON(200, "tres coooool")
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
