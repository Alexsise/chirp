package main

import (
	"log"

	"chirp/models"
	"chirp/routes"

	"github.com/gin-gonic/gin"
)

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

	// Swagger setup

	// Initialize routes
	routes.InitRoutes(r, db)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
