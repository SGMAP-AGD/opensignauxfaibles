package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

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
		Key:        []byte("this is my secret key"),
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

		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
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
		api.GET("/database/status", getDBStatus)
		api.GET("/mock/compact", mockCompact)
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
		api.GET("/lastMove", lastMove)
	}
	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func lastMove(c *gin.Context) {
	dbstatus := c.Keys["DBSTATUS"].(*mgo.Database)

	var lastMove struct {
		ID       AdminID `json:"id" bson:"_id"`
		LastMove int     `json:"last_move" bson:"last_move"`
	}

	dbstatus.C("Admin").Find(bson.M{"_id.type": "last_move", "_id.key": "last_move"}).One(&lastMove)

	c.JSON(200, lastMove.LastMove)
}

func mockCompact(c *gin.Context) {
	dbstatus := c.Keys["DBSTATUS"].(*mgo.Database)
	var status DBStatus
	dbstatus.C("Admin").Find(bson.M{"_id.key": "status", "_id.type": "status"}).One(&status)
	s := "Traitement du lot 1807"
	status.Status = &s
	dbstatus.C("Admin").Upsert(bson.M{"_id": status.ID}, status)

	var lastMove struct {
		ID       AdminID `json:"id" bson:"_id"`
		LastMove int     `json:"last_move" bson:"last_move"`
	}

	dbstatus.C("Admin").Find(bson.M{"_id.type": "last_move", "_id.key": "last_move"}).One(&lastMove)
	lastMove.LastMove++
	dbstatus.C("Admin").Upsert(bson.M{"_id": lastMove.ID}, lastMove)

	go func() {
		time.Sleep(30 * time.Second)
		mockFree(c)
	}()
}

func mockFree(c *gin.Context) {
	dbstatus := c.Keys["DBSTATUS"].(*mgo.Database)
	var status DBStatus
	dbstatus.C("Admin").Find(bson.M{"_id.key": "status", "_id.type": "status"}).One(&status)
	status.Status = nil
	dbstatus.C("Admin").Upsert(bson.M{"_id": status.ID}, status)

	var lastMove struct {
		ID       AdminID `json:"id" bson:"_id"`
		LastMove int     `json:"last_move" bson:"last_move"`
	}

	dbstatus.C("Admin").Find(bson.M{"_id.type": "last_move", "_id.key": "last_move"}).One(&lastMove)
	lastMove.LastMove++
	dbstatus.C("Admin").Upsert(bson.M{"_id": lastMove.ID}, lastMove)
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
