package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/Laevateinn17/travelohi-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "vincent ganteng"

func HandleRegister(c *fiber.Ctx) error {

	var data models.Payload

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("Error binding JSON: %s\n", err)
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": "Invalid JSON"})
	}

	database, _ := db.Connect()
	err := RegisterUser(database, &data.User, &data.UserAuth)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return nil
}

func Ping(c *fiber.Ctx) error {
	c.SendStatus(http.StatusOK)
	return nil
}

func HandleLogin(c *fiber.Ctx) error {
	var userAuth models.UserAuth
	if err := c.BodyParser(&userAuth); err != nil {
		fmt.Println("Error binding json")
	}

	database, _ := db.Connect()
	user, err := GetUserAuth(database, &userAuth)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(err.Error())
	}

	expiryTime := time.Now().Add(time.Hour * 24)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(expiryTime),
	})

	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.JSON(fiber.Map{
			"error": "could not login",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiryTime,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	cookie2 := fiber.Cookie{
		Name:    "session",
		Value:   "active",
		Expires: expiryTime,
	}

	c.Cookie(&cookie2)

	return c.JSON(fiber.Map{
		"messsage": "success",
	})
}

func GetUserData(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims := token.Claims

	id, err := claims.GetIssuer()

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	database, _ := db.Connect()
	var userAuth models.UserAuth
	err = database.Model(&models.UserAuth{}).Where("id = ?", id).First(&userAuth).Error

	if err == nil {
		var user models.User
		database.Model(&user).Where("id = ?", userAuth.UserID).First(&user)
		return c.JSON(user)
	}
	// fmt.Println(err.Error())
	return nil
}

func HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "deleted",
		Expires:  time.Now().Add(-(time.Hour * 2)), // Add negative time means it happens in the past :P
		HTTPOnly: true,
	})

	c.ClearCookie("session")
	fmt.Println("clearing cookies")
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func CreateOTPRequest(c *fiber.Ctx) error {
	var userAuth models.UserAuth

	if err := c.BodyParser(&userAuth); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed connecting to database")
	}
	fmt.Println(userAuth.Email)
	if len(userAuth.Email) <= 0 || !models.DoesEmailExist(database, userAuth.Email) {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf("supplied email is invalid")
	}

	otp, err := CreateOTP(database, &userAuth)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed sending otp code")
	}

	c.Status(http.StatusOK)

	return c.JSON(otp)
}

func ValidateOTPRequest(c *fiber.Ctx) error {
	var otp models.OTP

	if err := c.BodyParser(&otp); err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("error binding json")
	}

	database, err := db.Connect()

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed connecting to database")
	}

	err = ValidateOTP(database, &otp)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(err.Error())
	}

	c.Status(http.StatusOK)
	return nil
}
