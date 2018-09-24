package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

func debug(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batches := getBatches(db)
	for _, b := range batches {
		b.Open = true
		b.Draft = true
		b.save(db)
	}
}
