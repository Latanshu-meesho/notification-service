package repositories

import (
	"notification-service/models"

	"gorm.io/gorm"
)

func GetSMSDetails(requestID string, db *gorm.DB) (*models.SMSRequest, error) {
	var smsRequest models.SMSRequest
	if err := db.Where("id = ?", requestID).First(&smsRequest).Error; err != nil {
		return nil, err
	}
	return &smsRequest, nil
}
func UpdateSMSStatus(requestID, status, failureDetails string, db *gorm.DB) error {
	return db.Exec(`
		UPDATE sms_requests
		SET status = ?, failure_comments = ?
		WHERE id = ?
	`, status, failureDetails, requestID).Error
}
func SaveSMSRequest(db *gorm.DB, smsRequest *models.SMSRequest) error {
	return db.Create(smsRequest).Error
}
