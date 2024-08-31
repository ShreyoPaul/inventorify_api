package controllers

import (
	"context"
	"fmt"
	model "main/DB"
	"main/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string
	Email string
}

func StoreBill(c *gin.Context) {
	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.D{{Key: "email", Value: claims.Email}}
	var existingUser Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var bill model.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	// currTime, err := time.Now().Format()
	if err != nil { // Always check errors even if they should not happen.
		panic(err)
	}

	bill.U_id = existingUser.ID
	bill.U_email = existingUser.Email
	bill.CreatedAt = time.Now().Local()

	res, err := model.Bills.InsertOne(context.Background(), bill)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"msg": "Bill created!", "result": res})
}

func GetAllBills(c *gin.Context) {
	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.D{{Key: "email", Value: claims.Email}}
	var existingUser Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser.ID)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	findOption := bson.D{{"u_email", existingUser.Email}, {"u_id", existingUser.ID}}

	cursor, err := model.Bills.Find(context.TODO(), findOption)
	if err != nil {
		panic(err)
	}
	var bills []model.Bill
	if err = cursor.All(context.TODO(), &bills); err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"msg": "Bill fetched!", "result": bills, "name": existingUser.Name, "email": existingUser.Email})
}
