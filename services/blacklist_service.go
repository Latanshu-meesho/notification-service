package services

import (
	"notification-service/repositories"
	"notification-service/utils"

	"gorm.io/gorm"
)

func GetBlacklistedNumbers() ([]string, error) {
	return utils.GetBlacklistFromRedis()
}

func AddToBlacklist(phoneNumbers []string, db *gorm.DB) error {
	// Add to Redis
	err := utils.AddToRedisBlacklist(phoneNumbers)
	if err != nil {
		return err
	}

	// Add to DB
	return repositories.AddToBlacklist(phoneNumbers, db)
}

func RemoveFromBlacklist(phoneNumbers []string, db *gorm.DB) error {
	// Remove from Redis
	err := utils.RemoveFromRedisBlacklist(phoneNumbers)
	if err != nil {
		return err
	}

	// Remove from DB
	return repositories.RemoveFromBlacklist(phoneNumbers, db)
}
