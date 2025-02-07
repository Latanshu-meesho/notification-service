package controllers

import (
	"log"
	"net/http"
	"notification-service/repositories"
	"notification-service/services"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SMSController handles SMS-related endpoints
type SMSController struct {
	Producer *kafka.Producer
	DB       *gorm.DB
	Topic    string
}

// SendSMSHandler handles the POST /v1/sms/send endpoint
func (ctrl *SMSController) SendSMSHandler(c *gin.Context) {
	var smsRequest struct {
		PhoneNumber string `json:"phoneNumber"` // JSON keys must match exactly
		Message     string `json:"message"`
	}
	// Parse the JSON body
	if err := c.ShouldBindJSON(&smsRequest); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body. Ensure all required fields are provided.",
			},
		})
		return
	}
	log.Printf("The request came to the server: %v", smsRequest)
	// Call the SendSMS service
	requestID, err := services.SendSMS(
		smsRequest.PhoneNumber,
		smsRequest.Message,
		ctrl.Producer,
		ctrl.Topic,
		ctrl.DB,
	)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process SMS request. Please try again later.",
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"requestId": requestID,
			"comments":  "Successfully Sent",
		},
	})

}

// GetSMSDetails handles the GET /v1/sms/:requestId endpoint
func (ctrl *SMSController) GetSMSDetails(c *gin.Context) {
	requestID := c.Param("requestId")

	// Fetch SMS details from the database
	smsDetails, err := repositories.GetSMSDetails(requestID, ctrl.DB)
	if err != nil {
		log.Printf("Failed to retrieve SMS details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "request_id not found",
			},
		})
		return
	}

	// Return SMS details
	c.JSON(http.StatusOK, gin.H{"data": smsDetails})
}
