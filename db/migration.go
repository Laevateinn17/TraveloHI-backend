package db

import (
	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	db.AutoMigrate(
		&models.User{},
		&models.UserAuth{},
		&models.FlightTicket{},
		&models.Airline{},
		&models.Airplane{},
		&models.Airport{},
		&models.FlightSchedule{},
		&models.Passenger{},
		&models.OTP{},
	)

}
