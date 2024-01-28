package controllers

import (
	"time"

	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func GetAirplanes(db *gorm.DB) ([]*models.Airplane, error) {
	var result []*models.Airplane
	if err := db.Preload("Airline").Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func IsAirplaneAvailable(db *gorm.DB, airplaneID uint, startTime, endTime time.Time) (bool, error) {
	var airplane models.Airplane

	if err := db.Preload("FlightSchedules").First(&airplane, airplaneID).Error; err != nil {
		return false, err
	}

	for _, schedule := range airplane.FlightSchedules {
		if !(endTime.Before(schedule.DepartureTime) || startTime.After(schedule.ArrivalTime)) {
			return false, nil
		}
	}

	return true, nil
}
