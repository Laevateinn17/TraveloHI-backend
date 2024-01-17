package main

import (
	"fmt"
	"net/http"

	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/Laevateinn17/travelohi-backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var connString = "host=localhost port=5432 user=postgres password=admin dbname=travelohi sslmode=disable"

func handleRegister(c *gin.Context) {

	var data models.Payload

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Printf("Error binding JSON: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
		return
	}

	database, _ := db.Connect(connString)
	err := models.RegisterUser(database, &data.User, &data.UserAuth)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

}

func testing(c *gin.Context) {
	c.Header("A", "")
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func handleLogin(c *gin.Context) {
	var userAuth models.UserAuth

	if err := c.ShouldBindJSON(&userAuth); err != nil {
		fmt.Println("Error binding json la\n")
	}

	fmt.Println("useryaaaaa ", userAuth)

	database, _ := db.Connect(connString)
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

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/ping", testing)
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)

	database, err := db.Connect(connString)

	if err != nil {
		fmt.Println("Error while connecting to db")
	}

	db.Migrate(database)

	router.Run("0.0.0.0:8080")
}
