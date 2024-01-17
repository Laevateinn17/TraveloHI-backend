package models

import "gorm.io/gorm"

type FlightTicket struct {
	gorm.Model
	FlightID   string      `json:"flightID" gorm:"foreignKey:ID"`
	Passengers []Passenger `json:"passengers" gorm:"foreignKey:TicketID"`
}
