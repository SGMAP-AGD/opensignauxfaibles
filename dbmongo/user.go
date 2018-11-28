package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// login
type login struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	BrowserToken string `form:"browserToken" json:"browserToken"`
	CheckCode    string `form:"checkCode" json:"checkCode"`
}

// AdminUser object utilisateur mongodb
type AdminUser struct {
	ID             AdminID   `json:"_id" bson:"_id"`
	HashedPassword []byte    `json:"hashedPassword,omitempty" bson:"hashedPassword,omitempty"`
	HashedRecovery []byte    `json:"hashedRecovery,omitempty" bson:"hashedRecovery,omitempty"`
	TimeRecovery   time.Time `json:"timeRecovery" bson:"timeRecovery"`
	HashedCode     []byte    `json:"hashedCode,omitempty" bson:"hashedCode,omitempty"`
	TimeCode       time.Time `json:"timeCode,omitempty" bson:"timeCode,omitempty"`
	Cookies        []string  `json:"cookies" bson:"cookies"`
	Level          string    `json:"level" bson:"level"`
	FirstName      string    `json:"firstName" bson:"firstName"`
	LastName       string    `json:"lastName" bson:"lastName"`
	BrowserTokens  []string  `json:"browserTokens" bson:"browserTokens"`
}

// Browser mappe les informations contenues dans un browserToken
type Browser struct {
	Name      string    `json:"name" bson:"name"`
	IP        string    `json:"ip" bson:"ip"`
	Created   time.Time `json:"created" bson:"created"`
	UserEmail string    `json:"email" bson:"email"`
}

func (user AdminUser) save() error {
	err := db.DBStatus.C("Admin").Update(bson.M{"_id": user.ID}, user)
	return err
}

// type AdminLevel string
const levelAdmin = "admin"
const levelPowerUser = "powerUser"
const levelUser = "user"

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	email := claims["id"].(string)
	user, err := loadUser(email)
	if err != nil {
		c.JSON(500, "Erreur d'identification")
	}
	return &user
}

func loginUser(username string, password string, browserToken string) (AdminUser, error) {
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

func loginUserWithCredentials(username string, password string) (AdminUser, error) {
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

func loadUser(email string) (AdminUser, error) {
	var user AdminUser
	if err := db.DBStatus.C("Admin").Find(bson.M{"_id.type": "credential", "_id.key": email}).One(&user); err != nil {
		return AdminUser{}, err
	}
	return user, nil
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginVals login

	if err := c.ShouldBind(&loginVals); err != nil {
		fmt.Println(err)
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	password := loginVals.Password
	browserToken := loginVals.BrowserToken
	user, err := loginUser(email, password, browserToken)

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

func getCode() int {
	i, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return int(i.Int64())
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

func checkRecoveryEmailHandler(c *gin.Context) {
	email := c.Params.ByName("email")
	code := c.Params.ByName("code")

	err := checkRecoveryEmail(email, code)
	if err != nil {
		c.JSON(401, err.Error())
	} else {
		c.JSON(200, nil)
	}
}

func checkRecoveryEmail(email string, code string) error {
	user, err := loadUser(email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(user.HashedRecovery, []byte(code))
	if err != nil {
		return err
	}
	user.HashedRecovery = nil
	user.TimeRecovery = time.Time{}
	err = user.save()
	return err
}

func sendRecoveryEmail(email string) error {
	user, err := loadUser(email)
	if err == nil {
		recoveryCode := fmt.Sprintf("%06d", getCode())
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(recoveryCode), bcrypt.DefaultCost)
		if err == nil {
			user.HashedRecovery = hashedPassword
			user.TimeRecovery = time.Now()
			err = user.save()
			spew.Dump(err)
			if err == nil {
				fmt.Println(recoveryCode)
			}
		}
	} else {
		fmt.Println("error: " + err.Error())
	}
	return nil

	// smtpAddress := viper.GetString("smtpAddress")
	// smtpUser := viper.GetString("smtpUser")
	// smtpPassword := viper.GetString("smtpPass")

	// c, err := smtp.Dial("localhost:25")
	// if err != nil {
	// 	spew.Dump(err)
	// }
	// defer c.Close()
	// // Set the sender and recipient.
	// c.Mail("christophe@zbouboutchi.net")
	// c.Rcpt("christophe@zbouboutchi.net")
	// // Send the email body.
	// wc, err := c.Data()
	// if err != nil {
	// 	spew.Dump(err)
	// }
	// defer wc.Close()
	// buf := bytes.NewBufferString("This is the email body.")
	// if _, err = buf.WriteTo(wc); err != nil {
	// 	spew.Dump(err)
	// }
	// return err

}

func loginGetHandler(c *gin.Context) {
	var loginVals login

	if err := c.ShouldBind(&loginVals); err != nil {
		fmt.Println(err)
		c.JSON(401, "requête malformée")
		return
	}

	err := loginGet(loginVals)

	if err != nil {
		c.JSON(500, "Erreur lors de l'envoi du code de validation")
	}

}

func loginGet(login login) error {
	email := login.Email
	password := login.Password
	user, err := loginUserWithCredentials(email, password)

	if err == nil {
		checkCode := fmt.Sprintf("%06d", getCode())
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(checkCode), bcrypt.DefaultCost)
		if err == nil {
			user.HashedCode = hashedCode
			user.TimeCode = time.Now()
			err = user.save()
			if err == nil {
				fmt.Println(checkCode)
			}
		}
	} else {
		fmt.Println("error: " + err.Error())
	}
	return err
}

func loginCheckHandler(c *gin.Context) {
	var loginVals login
	c.ShouldBind(&loginVals)

	email := loginVals.Email
	password := loginVals.Password
	checkCode := loginVals.CheckCode

	err := loginCheck(email, password, checkCode)

	if err != nil {
		c.JSON(401, "Authentification Incorrecte")
	} else {
		c.JSON(200, "Welcome")
	}
}

func loginCheck(email string, password string, checkCode string) error {
	fmt.Println(email, password, checkCode)

	user, err := loginUserWithCredentials(email, password)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedCode, []byte(checkCode))
	if err != nil {
		return err
	}

	user.HashedCode = nil
	user.TimeCode = time.Time{}
	err = user.save()
	return err
}
