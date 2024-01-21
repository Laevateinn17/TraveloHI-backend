package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Laevateinn17/travelohi-backend/utilities"
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
