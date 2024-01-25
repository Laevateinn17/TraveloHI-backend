package controllers

import (
	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func GetAirlines(db *gorm.DB) ([]*models.Airline, error) {
	var result []*models.Airline
	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}


	return result, nil
}