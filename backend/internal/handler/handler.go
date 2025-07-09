package handler

import (
	"net/http"

	"github.com/IlfGauhnith/Hackathon-Bix-3T-Golang/internal/logger"
	"github.com/gin-gonic/gin"
)

// Health Check godoc
// @Summary      Check API health
// @Description  Check if the API is running
// @Produce      json
// @Success      200  {object}  dtoResponse.ItemResponse
// @Failure      500  {object}  dtoResponse.ErrorResponse "Internal server error"
// @Router       /items/scan/{ean13} [get]
func HealthHandler(c *gin.Context) {
	logger.Log.Info("HealthHandler")
	c.JSON(http.StatusOK, gin.H{"status": "API is running"})
}
