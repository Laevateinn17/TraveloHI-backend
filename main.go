package main

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/controllers"
	"github.com/Laevateinn17/travelohi-backend/db"
	"github.com/Laevateinn17/travelohi-backend/seed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

func main() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
	}))

	router.Get("/ping", controllers.Ping)
	router.Post("/register", controllers.HandleRegister)
	router.Post("/login", controllers.HandleLogin)
	router.Post("/email-login", controllers.HandleLoginByEmail)

	router.Get("/user", controllers.GetUserData)
	router.Post("/logout", controllers.HandleLogout)

	router.Post("/otp", controllers.CreateOTPRequest)
	router.Post("/check-otp", controllers.ValidateOTPRequest)

	router.Post("/security-question", controllers.HandleGetSecurityQuestion)
	router.Post("/security-answer", controllers.ValidateSecurityAnswer)

	router.Post("/change-password", controllers.HandleChangePassword)

	router.Get("/airports-data", controllers.HandleGetAirportsData)

	router.Get("/airlines-data", controllers.HandleGetAirlinesData)

	router.Get("/airplanes-data", controllers.HandleGetAirplanesData)

	router.Post("/add-flight", controllers.HandleAddFlightSchedule)

	router.Post("/available-planes", controllers.HandleGetAvailableAirplanes)

	database, err := db.Connect()

	if err != nil {
		fmt.Println("Error while connecting to db")
		return
	}

	db.Migrate(database)

	seed.Seed(database)

	router.Listen("0.0.0.0:8080")
}
