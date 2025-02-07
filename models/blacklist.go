package models

type BlacklistedNumbers struct {
	PhoneNumbers string `json:"phoneNumber" binding:"required"`
}
