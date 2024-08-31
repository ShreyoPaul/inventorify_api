package controllers

import (
	"context"
	"fmt"
	model "main/DB"
	"main/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddAttributesRequestBody struct {
	Attributes  []string
	InventoryId primitive.ObjectID
}

func GetAllInv(c *gin.Context) {
	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.M{"email": claims.Email}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser.Inv)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	c.JSON(200, gin.H{"success": "Fetched all inventory!", "inventories": existingUser.Inv, "name": existingUser.Name, "email": existingUser.Email})
}

func CreateInv(c *gin.Context) {
	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.D{{Key: "email", Value: claims.Email}}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser.Name)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var inventory model.Inv
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	inventory.ID = primitive.NewObjectID()

	update := bson.M{
		"$push": bson.M{"inv": inventory},
	}
	result, err := model.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(400, gin.H{"error": "Inventory creation failed!"})
		return
	}

	c.JSON(200, gin.H{"success": "home page", "result": result})
}

func AddAttributes(c *gin.Context) {
	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.D{{Key: "email", Value: claims.Email}}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	fmt.Println("test", existingUser.Name)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var addAttributesRequestBody AddAttributesRequestBody
	if err := c.ShouldBindJSON(&addAttributesRequestBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Print("Array:", addAttributesRequestBody)

	filterbyId := bson.M{"inv._id": addAttributesRequestBody.InventoryId}

	update := bson.M{
		"$push": bson.M{
			"inv.$.inv_attributes": bson.M{
				"$each": addAttributesRequestBody.Attributes,
			},
		},
	}

	result, err := model.Collection.UpdateOne(context.Background(), filterbyId, update)
	if err != nil {
		c.JSON(400, gin.H{"error": "Inventory creation failed! " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": result})
}
