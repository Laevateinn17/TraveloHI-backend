package models

import "gorm.io/gorm"

type FlightTicket struct {
	gorm.Model
	Passengers     []Passenger      `json:"passengers" gorm:"foreignKey:FlightTicketID"`
	FlightSchedule []FlightSchedule `gorm:"foreignKey:FlightTicketID"`
}
