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
	err := models.RegisterUser(database, &data.User, &data.UserAuth)

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

	fmt.Println("useryaaaaa ", userAuth)

	database, _ := db.Connect()
	user, err := models.GetUser(database, &userAuth)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
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
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}
	fmt.Println(time.Now().String())
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"messsage": "success",
	})
}
