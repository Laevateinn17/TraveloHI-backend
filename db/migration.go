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
		&models.Airplane{},
		&models.FlightSchedule{},
		&models.Airline{},
		&models.Airport{},
		&models.Passenger{},
		&models.OTP{},
	)

}
