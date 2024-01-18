package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Laevateinn17/travelohi-backend/utilities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName         string    `json:"firstName" gorm:"column:first_name"`
	MiddleName        string    `json:"middleName" gorm:"column:middle_name"`
	LastName          string    `json:"lastName" gorm:"column:last_name"`
	DateOfBirth       time.Time `json:"dateOfBirth"`
	Gender            string    `json:"gender"`
	IsBanned          bool      `json:"isBanned"`
	ProfilePictureURL string    `json:"profilePictureURL"`

	UserAuth UserAuth `json:"-" gorm:"foreignKey:UserID"`
}

type UserAuth struct {
	gorm.Model
	UserID           uint   `gorm:"foreignKey"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	SecurityQuestion string `json:"securityQuestion"`
	SecurityAnswer   string `json:"securityAnswer"`
}

type Payload struct {
	User     User     `json:"user"`
	UserAuth UserAuth `json:"userAuth"`
}

func DoesEmailExist(db *gorm.DB, email string) bool {
	var count int64
	db.Model(&UserAuth{}).Where("email = ?", email).Count(&count)

	return count > 0
}

func ValidateData(user *User, userAuth *UserAuth) bool {

	if len(user.FirstName) <= 5 || utilities.HasNumber(user.FirstName) || utilities.HasSymbol(user.FirstName) {
		fmt.Println("1")
		return false
	}

	if time.Now().Year()-user.DateOfBirth.Year() < 13 {
		fmt.Println("2")
		return false
	}

	if user.Gender != "M" && user.Gender != "F" {
		fmt.Println("3")
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+{}\\[\\]:;<>,.?~\\\\/-]+$`)

	if re.MatchString(userAuth.Password) || len(userAuth.Password) < 8 || len(userAuth.Password) > 30 {
		fmt.Println("4")
		return false
	}

	return true
}

func GetUser(db *gorm.DB, userAuth *UserAuth) (*User, error) {

	var result UserAuth
	db.Model(&UserAuth{}).Where("email = ?", userAuth.Email).First(&result)

	if result.ID == 0 {
		return nil, fmt.Errorf("authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(userAuth.Password)); err != nil {
		return nil, fmt.Errorf("invalid credential")
	}

	var user User
	db.Model(&User{}).Where("id = ?", result.UserID).First(&user)
	return &user, nil

}

func RegisterUser(db *gorm.DB, user *User, userAuth *UserAuth) error {

	if !ValidateData(user, userAuth) || DoesEmailExist(db, userAuth.Email) {
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
