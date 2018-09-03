package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
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
	r.Use(gin.Recovery())
	r.Use(DB())
	r.Use(Kanboard())
	// FIXME: configurer correctement CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST", "PUT", "HEAD", "DELETE")

	r.Use(cors.New(config))

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:         "OpenSignauxFaibles",
		Key:           []byte(viper.GetString("JWT_SECRET")),
		Timeout:       5 * time.Minute,
		MaxRefresh:    5 * time.Minute,
		Authenticator: authenticator,
		Authorizator:  authorizator,
		Unauthorized:  unauthorized,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	r.POST("/login", authMiddleware.LoginHandler)

	api := r.Group("api")
	api.Use(authMiddleware.MiddlewareFunc())

	{
		api.GET("/refreshToken", authMiddleware.RefreshHandler)
		api.GET("/purge", purge)
		api.GET("/kanboard/get/projects", listProjects)
		api.GET("/kanboard/get/tasks", getKBTasks)
		api.POST("/admin/batch", upsertBatch)
		api.GET("/admin/batch", listBatch)
		api.DELETE("/admin/batch", dropBatch)
		api.GET("/admin/files", adminFiles)
		api.POST("/admin/attach", attachFileBatch)
		api.GET("/admin/types", listTypes)
		api.GET("/admin/clone/:to", cloneDB)
		api.GET("/admin/features", adminFeature)
		api.GET("/admin/status", getDBStatus)
		api.GET("/batch/reset", resetBatch)
		api.GET("/batch/purge", purgeBatch)
		api.GET("/batch/process", processBatch)
		api.GET("/data/naf", getNAF)
		api.GET("/lastMove", lastMove)

		api.GET("/data/prediction/:batch/:algo/:page", predictionBrowse)
		api.GET("/debug/", debug)
		api.GET("/import/:batch", importBatch)
		api.GET("/compact/etablissement/:siret", compactEtablissement)
		api.GET("/compact/etablissement", compactEtablissement)
		api.GET("/compact/entreprise/:siren", compactEntreprise)
		api.GET("/compact/entreprise", compactEntreprise)
		api.GET("/reduce/:algo/:batch/:siret", reduce)
		api.GET("/reduce/:algo/:batch", reduce)
	}

	debugAPI := r.Group("debugAPI")
	debugAPI.Use(authMiddleware.MiddlewareFunc())
	{
		debugAPI.GET("/data/prediction/:batch/:algo/:page", predictionBrowse)
		debugAPI.GET("/debug/", debug)
		debugAPI.GET("/import/:batch", importBatch)
		debugAPI.GET("/compact/etablissement/:siret", compactEtablissement)
		debugAPI.GET("/compact/etablissement", compactEtablissement)
		debugAPI.GET("/compact/entreprise/:siren", compactEntreprise)
		debugAPI.GET("/compact/entreprise", compactEntreprise)
		debugAPI.GET("/reduce/:algo/:batch/:siret", reduce)
		debugAPI.GET("/reduce/:algo/:batch", reduce)
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

func purgeBatch(c *gin.Context) {}

func resetBatch(c *gin.Context) {}
