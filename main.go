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

func register(c *gin.Context) {
	

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

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	
	router.GET("/testing", testing)
	router.POST("/register", register)


	
	database, err := db.Connect(connString)
	
	if err != nil {
		fmt.Println("Error while connecting to db")
	}
	
	db.Migrate(database)




	router.Run("0.0.0.0:8080")
}