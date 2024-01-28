package models

import "gorm.io/gorm"

type Airline struct {
	gorm.Model

	Name      string     `json:"name"`
	Airplanes []Airplane `gorm:"foreignKey:AirlineID" json:"-"`
}
