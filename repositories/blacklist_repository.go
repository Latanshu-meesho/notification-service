package repositories

import (
	"log"

	"gorm.io/gorm"
)

func AddToBlacklist(phoneNumbers []string, db *gorm.DB) error {
	for _, number := range phoneNumbers {
		err := db.Exec("INSERT INTO blacklisted_numbers (phone_numbers) VALUES (?)", number).Error
		if err != nil {
			log.Printf("Failed to insert phone number %s: %v", number, err)
			return err
		}
	}
	return nil
}

func RemoveFromBlacklist(phoneNumbers []string, db *gorm.DB) error {
	for _, number := range phoneNumbers {
		err := db.Exec("DELETE FROM blacklisted_numbers WHERE phone_numbers = ?", number).Error
		if err != nil {
			return err
		}
	}
	return nil
}
