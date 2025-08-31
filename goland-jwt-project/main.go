package main

import (
	"context"
	"github.com/Somu/golang-jwt-project/routes"
	"log"
	"os"

	"github.com/Somu/golang-jwt-project/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"

	}
	router := gin.New()
	database.DBinstance()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		collection := database.OpenCollection(database.Client, "user")
		var user bson.M
		collection.FindOne(context.Background(), bson.M{"name": "test"}).Decode(&user)
		c.JSON(200, gin.H{
			"message": user,
		})

	})
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from API 2 telling access granted",
		})
	})

	router.Run(":" + port)

}
