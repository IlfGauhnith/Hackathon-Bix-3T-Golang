package main

import (
	"os"
	"time"

	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/internal/handler"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Hackathon Bix 3T Golang API
// @version 1.0
// @description This is the API documentation for Hackathon Bix 3T Golang API
// @BasePath /
func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080" // fallback
	}

	logger.Log.Info("Starting API server")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Store-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Browser can cache this config for 12 hours
	}))

	// Initialize routes from the router package
	initRoutes(router)
	logger.Log.Info("Routes successfully initialized.")

	router.Run("0.0.0.0:" + port)
	logger.Log.Infof("API server listening on port %s", port)
}

func initRoutes(router *gin.Engine) {
	// Support browser preflight OPTIONS requests explicitly
	// Ensures the browser receives the appropriate headers even on preflight requests
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	// Health endpoint
	router.GET("/health", handler.HealthHandler)

	docsGroup := router.Group("/docs")
	{
		// Serve Redoc HTML
		docsGroup.StaticFile("/openapi", "./docs/docs.html")

		// Serve raw YAML
		docsGroup.StaticFile("/openapi.yaml", "./docs/swagger.yaml")
	}

}
