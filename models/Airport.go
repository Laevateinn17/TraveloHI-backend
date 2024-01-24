package models

import "gorm.io/gorm"

type Airport struct {
	gorm.Model
	Name    string `json:"name"`
	Code    string `json:"code"`
	City    string `json:"city"`
	Country string `json:"country"`

	DepartureFlights []FlightSchedule `gorm:"foreignKey:DepartureAirportID"`
	ArrivalFlights   []FlightSchedule `gorm:"foreignKey:DestinationAirportID"`
}
