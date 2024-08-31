package main

import (
	"log"
	model "main/DB"
	"main/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://inventorify-iyct.vercel.app", "http://localhost:3000"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// r.Use(utils.CORSMiddleware())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	model.Init()
	routes.Auth(r)
	routes.Inventory(r)
	routes.Product(r)
	routes.Billing(r)
	r.Run(":8080")
}
