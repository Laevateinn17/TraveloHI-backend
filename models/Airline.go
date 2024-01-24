package models

import "gorm.io/gorm"

type Airline struct {
	gorm.Model

	Name      string     `json:"airline"`
	Airplanes []Airplane `gorm:"foreignKey:AirlineID" json:"-"`
}
