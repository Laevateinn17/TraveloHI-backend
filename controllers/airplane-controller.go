package controllers

import (
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