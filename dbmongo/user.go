package main

import "github.com/gin-gonic/gin"

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func authenticator(userID string, password string, c *gin.Context) (interface{}, bool) {
	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &User{
			UserName:  userID,
			LastName:  "Bo-Yi",
			FirstName: "Wu",
		}, true
	}

	return nil, false
}

func authorizator(user interface{}, c *gin.Context) bool {
	if v, ok := user.(string); ok && v == "admin" {
		return true
	}

	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
