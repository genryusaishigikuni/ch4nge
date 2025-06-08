// main.go
package main

import (
	"fmt"
	"github.com/genryusaishigikuni/ch4nge/config"
	"github.com/genryusaishigikuni/ch4nge/database"
	"github.com/genryusaishigikuni/ch4nge/handlers/seeder"
	"github.com/genryusaishigikuni/ch4nge/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set Gin mode based on environment
	if config.AppConfig.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	database.InitDB()

	if err := seeder.InitializeBasicData(); err != nil {
		log.Fatal("Failed to initialize basic data:", err)
	}
	// Setup Gin
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"environment": config.AppConfig.Server.Env,
		})
	})

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	address := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	fmt.Printf("Server starting on %s in %s mode\n", address, config.AppConfig.Server.Env)
	log.Fatal(r.Run(address))
}
