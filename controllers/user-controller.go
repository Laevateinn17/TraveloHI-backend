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
	fmt.Println(1)
	if result.ID == 0 {
		return nil, fmt.Errorf("authentication failed")
	}

	fmt.Println(result.Password)
	fmt.Println("password: ", userAuth.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("invalid credential")
	}

	fmt.Println(3)
	var user models.User
	db.Model(&models.User{}).Where("id = ?", result.UserID).First(&user)
	fmt.Println(4)
	return &user, nil

}

func GetUserAuthByEmail(db *gorm.DB, email string) (*models.UserAuth, error) {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", email).First(&result)
	if result.ID == 0 {
		return nil, fmt.Errorf("authentication failed")
	}

	return &result, nil
}

func GetUserAuth(db *gorm.DB, userAuth *models.UserAuth) (*models.UserAuth, error) {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		fmt.Println(2)
		return nil, fmt.Errorf("authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err != nil {
		fmt.Println(3)
		return nil, fmt.Errorf("invalid credential")
	}

	fmt.Println(4)
	return &result, nil

}

func GetSecurityQuestion(db *gorm.DB, userAuth *models.UserAuth) (string, error) {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		return "", fmt.Errorf("authentication failed")
	}

	return result.SecurityQuestion, nil
}

func GetSecurityAnswer(db *gorm.DB, userAuth *models.UserAuth) (string, error) {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		return "", fmt.Errorf("authentication failed")
	}

	return result.SecurityAnswer, nil
}

func ChangePassword(db *gorm.DB, userAuth *models.UserAuth) error {
	var result models.UserAuth
	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err == nil {
		return fmt.Errorf("new password is the same as old one")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("an unexpected error occurred")
	}

	db.Model(&models.UserAuth{}).Where("email = ?", userAuth.Email).Update("password", hash)
	return nil
}

func RegisterUser(db *gorm.DB, user *models.User, userAuth *models.UserAuth) error {

	if !models.ValidateData(user, userAuth) {
		return fmt.Errorf("one or more validation is violated")
	}

	if models.DoesEmailExist(db, userAuth.Email) {
		return fmt.Errorf("email is already used")
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
