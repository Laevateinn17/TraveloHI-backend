package controllers

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, userAuth *models.UserAuth) (*models.User, error) {

	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		return nil, fmt.Errorf("authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("invalid credential")
	}

	var user models.User
	db.Model(&models.User{}).Where("id = ?", result.UserID).First(&user)
	return &user, nil

}

func GetUserAuth(db *gorm.DB, userAuth *models.UserAuth) (*models.UserAuth, error) {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		return nil, fmt.Errorf("authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("invalid credential")
	}

	return &result, nil

}

func RegisterUser(db *gorm.DB, user *models.User, userAuth *models.UserAuth) error {

	if !models.ValidateData(user, userAuth) || models.DoesEmailExist(db, userAuth.Email) {
		return fmt.Errorf("one or more validation is violated")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("an unexpected error occurred")
	}

	userAuth.Password = string(hash)

	user.UserAuth = *userAuth
	result := db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("an unexpected error occurred")
	}
	return nil
}
