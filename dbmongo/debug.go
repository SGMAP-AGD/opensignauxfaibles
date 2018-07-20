package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func debug(c *gin.Context) {
	naf, err := loadNAF()
	spew.Dump(naf)
	spew.Dump(err)
}
