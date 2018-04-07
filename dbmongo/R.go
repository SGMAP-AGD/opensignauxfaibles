package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/senseyeio/roger"
)

func algo1(c *gin.Context) {
	Rscript, _ := ioutil.ReadFile("R/algo1.R")

	rClient, err := roger.NewRClient("127.0.0.1", 6311)
	if err != nil {
		fmt.Println("Failed to connect")
		return
	}

	value, err := rClient.Eval(string(Rscript))
	if err != nil {
		fmt.Println("Command failed: " + err.Error())
	} else {
		fmt.Println(value)
	}
}
