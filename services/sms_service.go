package services

import (
	"errors"
	"log"
	"notification-service/models"
	"notification-service/repositories"
	"notification-service/utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

var ErrFailedToGenerateID = errors.New("failed to generate unique request ID")

// SendSMS handles the SMS sending process
func SendSMS(phoneNumber string, message string, producer *kafka.Producer, topic string, db *gorm.DB) (string, error) {
	// Generate a unique request ID
	requestID := utils.GenerateUUID()
	if requestID == "" {
		return "", ErrFailedToGenerateID
	}

	// Save the SMS request in the database
	smsRequest := models.SMSRequest{
		ID:          requestID,
		PhoneNumber: phoneNumber,
		Message:     message,
		Status:      "PENDING",
	}
	err := repositories.SaveSMSRequest(db, &smsRequest)
	if err != nil {
		log.Printf("Failed to save SMS request to the mysql: %v", err)
		return "", err
	}

	// Publish the request ID to Kafka
	err = utils.PublishKafkaMessage(producer, topic, requestID)
	if err != nil {
		log.Printf("Failed to publish requestID to Kafka: %v", err)
		return "", err
	}

	log.Printf("SMS request created with ID: %s", requestID)
	done := make(chan bool)
	go utils.ConsumeMessages(db, done)

	<-done // Wait for message processing to complete
	return requestID, nil
}
