package controllers

import (
	"context"
	"fmt"
	model "main/DB"
	"main/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetProduct struct {
	// _id string
	Inv []model.Inv
}

type Product struct {
	ID         primitive.ObjectID
	Name       string
	Qty        uint64
	Amount     uint64
	Attributes []string
	CGST       float64
	SGST       float64
}

type Result struct {
	Products []Product
}

type PidRequest struct {
	Pid string `json:"pid"`
}

func GetAllProducts(c *gin.Context) {
	inventoryId := c.Param("inventory")
	id, err := primitive.ObjectIDFromHex(inventoryId)
	if err != nil {
		c.JSON(401, gin.H{"error": "Inventory Hex Error"})
		return
	}

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

	filterbyId := bson.M{"inv.id": id}
	projection := bson.M{"inv.$": 1}
	fmt.Println(id, inventoryId)
	var result GetProduct

	err = model.Collection.FindOne(context.Background(), filterbyId, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("result", result)

	c.JSON(200, gin.H{"user": existingUser.Name, "email": existingUser.Email, "result": result.Inv})
}

func CreateProduct(c *gin.Context) {
	inventoryId := c.Param("inventory")
	id, err := primitive.ObjectIDFromHex(inventoryId)
	if err != nil {
		c.JSON(401, gin.H{"error": "Inventory Hex Error"})
		return
	}

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
	// fmt.Println("test", existingUser.Inv)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	// var inventories []model.Inv
	// inventories = existingUser.Inv

	// var foundInv model.Inv

	// for _, inventory := range inventories {
	// 	if inventory.ID == id {
	// 		foundInv = inventory
	// 	}
	// }

	var products []model.Product
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// foundInv.Products = append(foundInv.Products, product)
	for i := range products {
		// Example: Print each item
		// fmt.Print(product, i)
		products[i].ID = primitive.NewObjectID()
		// fmt.Print(product.ID)
	}

	fmt.Println("product", products)
	filterbyId := bson.M{"inv.id": id}
	update := bson.M{
		"$push": bson.M{
			"inv.$.products": bson.M{
				"$each": products,
			},
		},
	}
	result, err := model.Collection.UpdateOne(context.Background(), filterbyId, update)
	if err != nil {
		c.JSON(400, gin.H{"error": "Inventory creation failed!" + err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "home page", "result": result})
}

func PipelineAllProducts(c *gin.Context) {
	// cookie, err := c.Cookie("token")
	// if err != nil {
	// 	c.JSON(401, gin.H{"error": "unauthorized! token not found", "cookie": cookie})
	// 	return
	// }

	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		fmt.Println("issue with auth token!")
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.M{"email": claims.Email}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	// fmt.Println("test", existingUser.Inv)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	fmt.Print("Testing")

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$inv"}},
		{{Key: "$unwind", Value: "$inv.products"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "products", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "id", Value: "$inv.products.id"},
				{Key: "name", Value: "$inv.products.name"},
				{Key: "qty", Value: "$inv.products.qty"},
				{Key: "amount", Value: "$inv.products.amount"},
				{Key: "cgst", Value: "$inv.products.cgst"},
				{Key: "sgst", Value: "$inv.products.sgst"},
			}}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "products", Value: 1},
		}}},
	}

	// var result []PipelineProduct

	cursor, err := model.Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var result Result
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("Products: %+v\n", result.Products)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	cursor.Close(context.TODO())

	c.JSON(200, gin.H{"user": existingUser.Name, "email": existingUser.Email, "result": result.Products})
}

func UpdateProduct(c *gin.Context) {
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
	// fmt.Println("test", existingUser.Inv)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	// collection := client.Database("your_database_name").Collection("your_collection_name")

	// updates := []UpdateDetails{
	// 	{ProductID: "66b27b7d20be3f5dd444f667", Qty: 5, Amount: 50},
	// 	{ProductID: "66b27b7d20be3f5dd444f668", Qty: 3, Amount: 5},
	// }

	var updates []model.Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// const res  *mongo.UpdateResult

	for _, update := range updates {
		fmt.Print(update)
		var id primitive.ObjectID = update.ID
		fmt.Print("ID--------->", update.ID)
		filter := bson.M{"inv.products.id": update.ID}
		update := bson.M{
			"$inc": bson.M{
				"inv.$.products.$[i].qty": update.Qty,
			},
			"$set": bson.M{
				"inv.$.products.$[i].amount": update.Amount,
				"inv.$.products.$[i].cgst":   update.CGST,
				"inv.$.products.$[i].sgst":   update.SGST,
			},
		}
		arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"i.id": id},
			},
		})
		fmt.Print("\nArray--------->", arrayFilters)

		// var updatedDoc bson.M
		result, err := model.Collection.UpdateOne(context.Background(), filter, update, arrayFilters)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Print(result)
		c.JSON(200, gin.H{"user": existingUser.Name, "email": existingUser.Email, "result": result})
	}
}

func BillingProduct(c *gin.Context) {
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
	fmt.Println("test", existingUser)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var updates []model.Product
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, update := range updates {
		fmt.Print(update)
		var id primitive.ObjectID = update.ID

		filter := bson.M{"inv.products.id": update.ID}
		update := bson.M{
			"$inc": bson.M{
				"inv.$[i].products.$.qty": -update.Qty,
			},
		}
		arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"i.products.id": id},
			},
		})
		// var updatedDoc bson.M
		result, err := model.Collection.UpdateOne(context.Background(), filter, update, arrayFilters)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Print(result)
		c.JSON(200, gin.H{"user": existingUser.Name, "email": existingUser.Email, "result": result})
	}
}

func DeleteProduct(c *gin.Context) {
	inventoryId := c.Param("inventory")
	id, err := primitive.ObjectIDFromHex(inventoryId)
	if err != nil {
		c.JSON(401, gin.H{"error": "Inventory Hex Error"})
		return
	}

	token := c.GetHeader("authorization")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized! Parsing failed!" + err.Error()})
		return
	}
	// fmt.Printf("Token claims added: %+v\n", claims.Email)

	filter := bson.M{"email": claims.Email}
	var existingUser model.Users
	model.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
	// fmt.Println("test", existingUser.Inv)
	if existingUser.Name == "" || existingUser.Email == "" {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var invs []model.Inv = existingUser.Inv
	// fmt.Print(invs)
	var pid PidRequest
	if err := c.ShouldBindJSON(&pid); err != nil {
		fmt.Print("ERROR in Pid\n")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// fmt.Print(pid)

	p, err := primitive.ObjectIDFromHex(pid.Pid)
	if err != nil {
		c.JSON(401, gin.H{"error": "Product Hex Error"})
		return
	}
	// fmt.Print(p)
	var newInv []model.Product
	for _, _inv := range invs {
		if _inv.ID == id {
			for _, product := range _inv.Products {
				if product.ID != p {
					newInv = append(newInv, product)
				}
			}
			fmt.Print("newInv:\n", newInv)

			_inv.Products = newInv
			fmt.Print("_inv.Products:\n", _inv.Products)
			break
		}
	}
	filterbyId := bson.M{"inv.id": id}
	update := bson.M{
		"$set": bson.M{
			"inv.$.products": newInv,
		},
	}
	result, err := model.Collection.UpdateOne(context.Background(), filterbyId, update)
	if err != nil {
		c.JSON(400, gin.H{"error": "Inventory creation failed!" + err.Error()})
		return
	}

	fmt.Print(result)
	c.JSON(200, gin.H{"user": existingUser.Name, "email": existingUser.Email, "result": result})
}
