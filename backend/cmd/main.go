// @title        Hackathon Bix 3T Golang API
// @version      1.0
// @description  This is the API documentation for Hackathon Bix 3T Golang API
// @host         localhost:8080
// @BasePath     /
package main

import (
	_ "net/http/pprof"

	"log"
	"net/http"
	"time"

	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/config"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/handler"
	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment-based configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	port := cfg.Port
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
	initRoutes(router, cfg)
	logger.Log.Info("Routes successfully initialized.")

	// Enable pprof for profiling
	go http.ListenAndServe("0.0.0.0:6060", nil) // pprof on :6060
	logger.Log.Info("pprof listening on port 6060")

	// Start HTTP server
	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func initRoutes(router *gin.Engine, cfg *config.Config) {
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

	// Upload endpoint for parallel processing
	router.POST("/upload", handler.UploadHandler(cfg))

	// Upload endpoint for sequential processing
	router.POST("/upload-seq", handler.UploadHandlerSequential(cfg))
}
