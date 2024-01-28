package controllers

import (
	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func GetAirports(db *gorm.DB) ([]*models.Airport, error) {
	var result []*models.Airport
	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
