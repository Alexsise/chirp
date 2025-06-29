package main

import (
	"log"

	"chirp/models"
	"chirp/routes"

	"github.com/gin-gonic/gin"

	_ "chirp/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer <your_token>

// @title Chirp API
// @version 1.0
// @description This is a Reddit-like REST API.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	db := models.InitDB()
	defer func() {
		dbSQL, err := db.DB()
		if err != nil {
			log.Fatal("Failed to close database connection:", err)
		}
		dbSQL.Close()
	}()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Swagger setup

	// Initialize routes
	routes.InitRoutes(r, db)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
