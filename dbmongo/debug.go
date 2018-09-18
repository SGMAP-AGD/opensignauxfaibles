package main

import (
	"github.com/gin-gonic/gin"
)

type debugType struct {
	Test int `json:"test" bson:"test"`
}

func debug(c *gin.Context) {
	batchID, err := nextBatchID("1812")
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, batchID)
	}
}
