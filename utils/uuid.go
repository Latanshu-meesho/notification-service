package utils

import (
	"log"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID.
func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return ""
	}
	return id.String()
}
