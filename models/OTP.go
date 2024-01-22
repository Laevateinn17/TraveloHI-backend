package models

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model 
	Email string `json:"email"`
	Code string `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}