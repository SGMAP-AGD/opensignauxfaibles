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
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (interface{}, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return &User{
					UserName:  userId,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, true
			}

			return nil, false
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			if v, ok := user.(string); ok && v == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	r.POST("/login", authMiddleware.LoginHandler)

	api := r.Group("api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/refresh_token", authMiddleware.RefreshHandler)

		api.GET("/purge", purge)
		api.GET("/kanboard/get/projects", listProjects)
		api.GET("/kanboard/get/tasks", getKBTasks)

		api.POST("/admin/batch", updateBatch)
		api.GET("/admin/batch", listBatch)
		api.DELETE("/admin/batch/:batchKey", dropBatch)

		api.GET("/admin/files", adminFiles)
		api.POST("/admin/attach", attachFileBatch)
		api.GET("/admin/types", listTypes)
		api.GET("/admin/clone/:to", cloneDB)
		api.GET("/admin/features", adminFeature)
		api.GET("/database/status", getDBStatus)
		api.GET("/database/lock", lockDBStatus)
		api.GET("/database/unlock", unlockDBStatus)
		api.GET("/import/:batch", importBatch)

		api.GET("/compact/etablissement/:siret", compactEtablissement)
		api.GET("/compact/etablissement", compactEtablissement)
		api.GET("/compact/entreprise/:siren", compactEntreprise)
		api.GET("/compact/entreprise", compactEntreprise)
		api.GET("/reduce/:algo/:batch/:siret", reduce)
		api.GET("/reduce/:algo/:batch", reduce)

		api.POST("/R/algo1", algo1)

		api.GET("/data/prediction/:batch/:algo/:page", predictionBrowse)
		api.GET("/data/naf", getNAF)
		api.GET("/debug/", debug)
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
