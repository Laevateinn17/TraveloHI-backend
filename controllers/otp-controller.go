package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)



func CreateOTP(db *gorm.DB, userAuth *models.UserAuth) (*models.OTP, error) {
	var result models.OTP

	db.Model(&models.OTP{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID != 0 {
		DeleteOTP(db, userAuth)
	}

	var otp models.OTP
	otp.Email = userAuth.Email
	otp.ExpiresAt = time.Now().Add(time.Minute * 5)

	for {
		otp.Code = fmt.Sprintf("%06d", rand.Intn(1000000))
		result.ID = 0
		db.Model(&models.OTP{}).Where("code = ?", otp.Code).First(&result)

		if result.ID == 0 {
			break
		}
	}

	db.Create(&otp)

	return &otp, nil
}

func ValidateOTP(db *gorm.DB, otp *models.OTP) error {
	var result models.OTP
	db.Model(&models.OTP{}).Where("code = ?", otp.Code).First(&result)

	if otp.ID == 0 || otp.Email != result.Email{
		return fmt.Errorf("invalid otp code")
	}

	if otp.ExpiresAt.Unix() < time.Now().Unix() {
		return fmt.Errorf("otp code has expired")
	}

	return nil
}

func DeleteOTP(db *gorm.DB, userAuth *models.UserAuth) {
	db.Delete(&userAuth)
}

