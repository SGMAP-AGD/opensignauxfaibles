package main

import (
	"bytes"
	"math/rand"
	"net/smtp"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// login
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AdminUser object utilisateur mongodb
type AdminUser struct {
	ID             AdminID `json:"_id" bson:"_id"`
	HashedPassword []byte  `json:"hashedPassword" bson:"hashedPassword"`
	Level          string  `json:"level" bson:"level"`
	FirstName      string  `json:"firstName" bson:"firstName"`
	LastName       string  `json:"lastName" bson:"lastName"`
}

// type AdminLevel string

const levelAdmin = "admin"
const levelPowerUser = "powerUser"
const levelUser = "user"

func loadUser(username string, password string) (AdminUser, error) {
	var user AdminUser
	if err := db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": username}).One(&user); err != nil {
		return AdminUser{}, err
	}
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == nil {
		return user, nil
	}
	return AdminUser{}, err

}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	user, err := loadUser(userID, password)

	if err == nil {
		return user, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*AdminUser); ok && v.ID.Key == "admin" {
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

func payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(AdminUser); ok {
		return jwt.MapClaims{
			identityKey: v.ID.Key,
		}
	}
	return jwt.MapClaims{}
}

func hashPassword(c *gin.Context) {
	password := c.Params.ByName("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, string(hashedPassword))
	}
}

func getRecoveryCode() int {
	rand.Seed(time.Now().UTC().UnixNano())
	recorveryCode := int((rand.Float64() * 1000000))
	return recorveryCode
}
func sendRecoveryEmailHandler(c *gin.Context) {
	address := c.Params.ByName("email")
	err := sendRecoveryEmail(address)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, nil)
	}

}

func sendRecoveryEmail(address string) error {
	// smtpAddress := viper.GetString("smtpAddress")
	// smtpUser := viper.GetString("smtpUser")
	// smtpPassword := viper.GetString("smtpPass")

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		spew.Dump(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail("christophe@zbouboutchi.net")
	c.Rcpt("christophe@zbouboutchi.net")
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		spew.Dump(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString("This is the email body.")
	if _, err = buf.WriteTo(wc); err != nil {
		spew.Dump(err)
	}
	return err
}
