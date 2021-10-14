package models

import (
	"gorm.io/gorm"
)

type Base struct {
	gorm.Model
	ID string `gorm:"type:uuid;primary_key;"`
}
