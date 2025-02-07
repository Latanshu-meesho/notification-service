package controllers

import (
	"log"
	"net/http"
	"notification-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlacklistController struct {
	DB *gorm.DB
}

// GetBlacklistedNumbers handles the GET /v1/blacklist endpoint
func (ctrl *BlacklistController) GetBlacklistedNumbers(c *gin.Context) {
	blacklist, err := services.GetBlacklistedNumbers()
	if err != nil {
		log.Printf("Error fetching blacklisted numbers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to retrieve blacklisted numbers.",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": blacklist})
}

// AddToBlacklist handles the POST /v1/blacklist endpoint
func (ctrl *BlacklistController) AddToBlacklist(c *gin.Context) {
	var req struct {
		PhoneNumbers []string `json:"phoneNumbers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request payload. Ensure phoneNumbers are provided.",
			},
		})
		return
	}
	err := services.AddToBlacklist(req.PhoneNumbers, ctrl.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Not able to blacklist",
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"data": "Successfully blacklisted"})
}

// RemoveFromBlacklist handles the DELETE /v1/blacklist endpoint
func (ctrl *BlacklistController) RemoveFromBlacklist(c *gin.Context) {
	var req struct {
		PhoneNumbers []string `json:"phoneNumbers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request payload. Ensure phoneNumbers are provided.",
			},
		})
		return
	}
	err := services.RemoveFromBlacklist(req.PhoneNumbers, ctrl.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Successfully whitelisted"})
}
