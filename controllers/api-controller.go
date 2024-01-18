package controllers

import (
	"fmt"
	"net/http"

	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/Laevateinn17/travelohi-backend/models"
	"github.com/gin-gonic/gin"
)

func HandleRegister(c *gin.Context) {

	var data models.Payload

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Printf("Error binding JSON: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
		return
	}

	database, _ := db.Connect()
	err := models.RegisterUser(database, &data.User, &data.UserAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

}

func Ping(c *gin.Context) {
	c.Header("A", "")
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func HandleLogin(c *gin.Context) {
	var userAuth models.UserAuth

	if err := c.ShouldBindJSON(&userAuth); err != nil {
		fmt.Println("Error binding json")
	}

	fmt.Println("useryaaaaa ", userAuth)

	database, _ := db.Connect()
	user, err := models.GetUser(database, &userAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	userResponse := models.UserResponse{
		Model:             user.Model,
		FirstName:         user.FirstName,
		MiddleName:        user.MiddleName,
		LastName:          user.LastName,
		DateOfBirth:       user.DateOfBirth,
		Gender:            user.Gender,
		IsBanned:          user.IsBanned,
		ProfilePictureURL: user.ProfilePictureURL,
	}

	c.IndentedJSON(http.StatusOK, userResponse)
}