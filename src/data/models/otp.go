package models

type OTP struct {
	Base
	Code   string
	Used   bool
	UserID string `gorm:"type:uuid;"`
	User   User
}
