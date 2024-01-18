package main

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/controllers"
	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)



func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/ping", controllers.Ping)
	router.POST("/register", controllers.HandleRegister)
	router.POST("/login", controllers.HandleLogin)


	database, err := db.Connect()
	if err != nil {
		fmt.Println("Error while connecting to db")
		return
	}


	db.Migrate(database)

	router.Run("0.0.0.0:8080")
}
