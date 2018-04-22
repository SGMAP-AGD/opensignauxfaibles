package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// User identification d'utilisateur
type User struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Permission []string      `json:"permission" bson:"permission"`
	Nom        string        `json:"nom" bson:"nom"`
	Prenom     string        `json:"prenom" bson:"prenom"`
	Contact    string        `json:"telephone" bson:"telephone"`
	Region     bson.ObjectId `json:"region_id" bson:"region_id"`
	jwt.StandardClaims
}

// Auth Informations de login
type Auth struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// JWTRequest Requête pour tester le token (mode POC)
type JWTRequest struct {
	Token string `json:"token" bson:"token"`
}

func createTokenString(user User) string {
	jwtSecret := viper.GetString("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &user)
	tokenstring, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}

func auth(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)

	var auth Auth
	var user User
	c.BindJSON(&auth)

	db.C("user").Find(bson.M{"_id": auth.Username, "password": auth.Password}).One(&user)

	if user.ID != "" {
		token := JWTRequest{
			Token: createTokenString(user),
		}
		c.JSON(200, token)
	} else {
		c.JSON(401, "Not authenticated")
	}

}

func readJWT(c *gin.Context) {
	jwtSecret := viper.GetString("JWT_SECRET")

	var request JWTRequest
	c.BindJSON(&request)

	clearToken, _ := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	c.JSON(200, clearToken)
}
