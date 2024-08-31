package model

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection
var Bills *mongo.Collection

func Init() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	Client = client

	log.Print("DB connected!")

	database := client.Database("GO_AUTH_1")
	database.Collection("Users")
	database.Collection("Bills")
	Collection = database.Collection("Users")
	Bills = database.Collection("Bills")
	log.Print(database.Collection("Users").Name())

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

}
