package controllers

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func AddFlightSchedule(db *gorm.DB, flightSchedule *models.FlightSchedule) error {
	var result models.FlightSchedule

	err := db.Model(&models.FlightSchedule{}).Where("flight_number = ?", flightSchedule.FlightNumber).First(&result).Error

	if result.ID != 0 {
		return fmt.Errorf("duplicate flight number")
	}

	if err = db.Create(&flightSchedule).Error; err != nil {
		return err
	}

	return nil

}
