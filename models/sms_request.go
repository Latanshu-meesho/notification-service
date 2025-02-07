package models

import (
	"time"
)

type SMSRequest struct {
	ID              string    `json:"id" gorm:"primaryKey"` // Primary Key
	PhoneNumber     string    `json:"phone_number" gorm:"not null"`
	Message         string    `json:"message" gorm:"not null"`
	Status          string    `json:"status"` // Status can be PENDING, SENT, FAILED, etc.
	FailureCode     string    `json:"failure_code"`
	FailureComments string    `json:"failure_comments"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
type DeliveredSMSRequest struct {
	DeliveryChannel string `json:"deliverychannel"`
	Channels        struct {
		SMS struct {
			Text string `json:"text"`
		} `json:"sms"`
	} `json:"channels"`
	Destination []struct {
		Msisdn        []string `json:"msisdn"`
		CorrelationID string   `json:"correlationId"`
	} `json:"destination"`
}
