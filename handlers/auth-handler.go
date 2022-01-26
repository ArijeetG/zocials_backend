package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"socials.com/helpers"
)

type Login struct {
	User     string
	Password string
}

type signUp struct {
	User     string
	Password string
	Email    string
}

type signUpPayload struct {
	uuid     string
	user     string
	password string
	email    string
	iat      int
}

type findUser struct {
	email    string
	iat      int
	password string
	user     string
	uuid     string
}

type Claims struct {
	uuid string
	jwt.StandardClaims
}

func AuthHandler(c *gin.Context) int {
	// client := helpers.ConnectToDatabase()
	// userCol := client.Database("")
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 0
	}
	if json.User == "" || json.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_parameter",
		})
		return 0
	}
	c.JSON(http.StatusOK, gin.H{"status": "correct_parameters"})
	return 1

}

func SignUp(c *gin.Context) int {
	client := helpers.ConnectToDatabase()
	usersCol := client.Database("master").Collection("Users")
	var json signUp
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 0
	}

	if json.User == "" || json.Email == "" || json.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_parameter",
		})
		return 0
	}

	isUserPresent, _ := usersCol.CountDocuments(context.TODO(), bson.D{{"user", json.User}})
	if isUserPresent >= 1 {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "user_present",
		})
		return 0
	}

	var payload signUpPayload
	uuid := uuid.New().String()
	bytes, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 14)
	payload.email = json.Email
	payload.uuid = uuid
	payload.user = json.User
	payload.password = string(bytes)
	payload.iat = int(time.Now().UnixMilli())

	doc := bson.D{{"uuid", payload.uuid}, {"email", payload.email}, {"user", payload.user}, {"password", payload.password}, {"iat", payload.iat}}

	_, err := usersCol.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_server_error",
		})

		return 0
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"status": "user_added",
		})
	}

	return 1
}

func Signin(c *gin.Context) int {
	client := helpers.ConnectToDatabase()
	usersCol := client.Database("master").Collection("Users")
	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}

	if json.Password == "" || json.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing_parameters",
		})
		return 0
	}

	isUserPresent := usersCol.FindOne(context.TODO(), bson.D{{"user", json.User}})
	if isUserPresent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_not_found",
		})
		return 0
	}
	var isUserPresentResult bson.M
	decodeError := isUserPresent.Decode(&isUserPresentResult)

	if decodeError != nil {
		// panic(decodeError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return 0
	}

	err := bcrypt.CompareHashAndPassword([]byte(isUserPresentResult["password"].(string)), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_password/email",
		})
	}

	//creating jwt-token
	var jwtKey = []byte("my-secret-key")
	claims := &Claims{
		uuid: isUserPresentResult["uuid"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(isUserPresentResult["iat"].(int64) + 48*3600*1000),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server error",
		})
		return 0
	}

	c.JSON(http.StatusAccepted, gin.H{
		"token": tokenString,
	})

	return 1
}
