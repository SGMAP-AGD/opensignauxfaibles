package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// Browser mappe les informations contenues dans un browserToken
type Browser struct {
	Name    string    `json:"name" bson:"name"`
	IP      string    `json:"ip" bson:"ip"`
	Created time.Time `json:"created" bson:"created"`
	Email   string    `json:"email" bson:"email"`
}

// BrowserToken embarque un token
type BrowserToken struct {
	BrowserToken string `json:"browserToken" bson:"browserToken"`
}

func forgeBrowserToken(browser Browser) (BrowserToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":    browser.Name,
		"ip":      browser.IP,
		"created": browser.Created,
		"email":   browser.Email,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("jwtSecret")))
	fmt.Println(err)
	return BrowserToken{
		BrowserToken: tokenString,
	}, err
}

func readBrowserToken(tokenString string) (Browser, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwtSecret")), nil
	})

	if err != nil {
		return Browser{}, err
	}
	if token.Valid {
		var browser Browser
		browser.Created, err = time.Parse("2006-01-02T15:04:05-0700", token.Claims.(jwt.MapClaims)["created"].(string))
		browser.Email = token.Claims.(jwt.MapClaims)["email"].(string)
		browser.IP = token.Claims.(jwt.MapClaims)["ip"].(string)
		browser.Email = token.Claims.(jwt.MapClaims)["email"].(string)
		return browser, nil
	}
	return Browser{}, nil
}
