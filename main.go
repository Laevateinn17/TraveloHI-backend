package main

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/controllers"
	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

func main() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.Get("/ping", controllers.Ping)
	router.Post("/register", controllers.HandleRegister)
	router.Post("/login", controllers.HandleLogin)

	database, err := db.Connect()
	if err != nil {
		fmt.Println("Error while connecting to db")
		return
	}

	db.Migrate(database)

	router.Listen("0.0.0.0:8080")
}
