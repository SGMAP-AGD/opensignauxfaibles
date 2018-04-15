package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	// Lancer Rserve en background
	go r()

	r := gin.Default()
	r.Use(DB())

	// FIXME: configurer correctement CORS
	r.Use(cors.Default())

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	v1 := r.Group("api/v1")
	{
		v1.GET("/purge", purge)
		v1.GET("/import", importData)
		v1.GET("/reduceEtablissement/:siret", reduceEtablissement)
		v1.GET("/reduceEtablissement", reduceEtablissements)
		v1.GET("/reduce/:siret", reduce)
		v1.GET("/reduce", reduceAll)
		v1.GET("/etablissement/:siret", browseEtablissement)
		v1.GET("/orig/:siret", browseOrig)
		v1.GET("/debug/:urssaf", debug)
		v1.GET("/importAP", importAP)
		v1.GET("/importDebit", importDebit)
		v1.GET("/importAltares", importAltares)
		v1.GET("/importEffectif", importEffectif)
		v1.POST("/R/algo1", algo1)
		v1.GET("/listFiles", listFiles)
		v1.GET("/data/debit/:siret", dataDebit)
	}

	r.Run(":3000")
}
