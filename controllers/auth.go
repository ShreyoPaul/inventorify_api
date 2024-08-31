package controllers

import (
	"context"
	"fmt"
	model "main/DB"
	"main/utils"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginBody struct {
	Pass  string
	Email string
}

var JwtKey = []byte("ilovegolang")

func Login(c *gin.Context) {
	var user model.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{{Key: "email", Value: user.Email}}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	claims := &model.Claims{
		Email: existingUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server error!"})
		return
	}
	fmt.Printf("Token claims added: %+v\n", tokenStr)

	c.SetCookie("token", tokenStr, 3600*7, "/", "localhost", true, false)
	c.JSON(201, gin.H{"msg": "Login successful!", "token": tokenStr})
}

func Signup(c *gin.Context) {
	var user model.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser model.Users
	filter := bson.D{{Key: "email", Value: user.Email}}
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser.Name)
	if existingUser.Name != "" || existingUser.Email != "" {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		c.JSON(500, gin.H{"error": "Could not generate password hash"})
		return
	}

	_, err := model.Collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
	}

	claims := &model.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server error!"})
		return
	}
	fmt.Printf("Token claims added: %+v\n", tokenStr)

	c.SetCookie("token", tokenStr, 3600*7, "/", "", false, false)
	c.JSON(200, gin.H{"msg": "Signup successful!", "token": tokenStr})
}
