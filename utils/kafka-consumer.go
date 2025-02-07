package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification-service/models"
	"notification-service/repositories"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func ConsumeMessages(db *gorm.DB, done chan bool) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notification.send_sms",
		GroupID: "notification-consumer-group",
	})

	defer reader.Close()
	ctx := context.Background()

	requestID, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message: %v", err)
		done <- false
		return
	}

	log.Printf("RequestID received to consumer: %s", string(requestID.Value))
	ProcessMessage(string(requestID.Value), db)
	done <- true
}

func ProcessMessage(requestID string, db *gorm.DB) {
	smsDetails, err := repositories.GetSMSDetails(requestID, db)
	if err != nil {
		log.Printf("Failed to fetch SMS details from mysql: %v", err)
		return
	}

	isBlacklisted, err := IsBlacklisted(smsDetails.PhoneNumber)
	if err != nil {
		log.Printf("Failed to check blacklist via Redis: %v", err)
		return
	}

	if isBlacklisted {
		log.Printf("Phone number %s is blacklisted", smsDetails.PhoneNumber)
		err = repositories.UpdateSMSStatus(requestID, "FAILED", "Phone number is blacklisted!", db)
		if err != nil {
			log.Printf("Failed to update SMS status: %v", err)
			return
		}
		return
	}
	err = DeliverSMS(smsDetails.PhoneNumber, smsDetails.Message, string(requestID))
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return
	}
	log.Printf("Sending SMS to %s: %s", smsDetails.PhoneNumber, smsDetails.Message)

	err = repositories.UpdateSMSStatus(requestID, "SENT", "", db)
	if err != nil {
		log.Printf("Failed to update SMS status: %v", err)
		return
	}

	doc := SMSDocument{
		ID:          smsDetails.ID,
		PhoneNumber: smsDetails.PhoneNumber,
		Message:     smsDetails.Message,
		CreatedAt:   smsDetails.CreatedAt,
		UpdatedAt:   smsDetails.UpdatedAt,
	}

	if err := IndexSMS(doc); err != nil {
		log.Printf("Failed to index SMS in Elasticsearch: %v", err)
		return
	}

	log.Printf("SMS processing completed for RequestID: %s", requestID)
}

func DeliverSMS(phoneNumber, messageText, requestID string) error {
	apiURL := "https://api.imiconnect.in/resources/v1/messaging"
	apiKey := "<API-KEY>"

	deliveredsmsRequest := models.DeliveredSMSRequest{
		DeliveryChannel: "sms",
		Destination: []struct {
			Msisdn        []string `json:"msisdn"`
			CorrelationID string   `json:"correlationId"`
		}{
			{
				Msisdn:        []string{phoneNumber},
				CorrelationID: requestID,
			},
		},
	}
	deliveredsmsRequest.Channels.SMS.Text = messageText

	requestBody, err := json.Marshal([]models.DeliveredSMSRequest{deliveredsmsRequest})
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS API returned status: %d", resp.StatusCode)
	}

	return nil
}
