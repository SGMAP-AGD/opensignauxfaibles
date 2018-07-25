package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

type debugType struct {
	Test int `json:"test" bson:"test"`
}

func debug(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)
	a := 1
	t1 := debugType{
		Test: a,
	}
	t2 := debugType{}
	db.C("test").Insert(t1)
	db.C("test").Insert(t2)
}
