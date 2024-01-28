package models

import "gorm.io/gorm"

type FlightTicket struct {
	gorm.Model
	Passengers      []Passenger      `json:"passengers" gorm:"foreignKey:FlightTicketID"`
	FlightSchedules []FlightSchedule `json:"flightSchedule" gorm:"many2many:flight_schedules_tickets"`
	Price           uint64           `json:"price"`
}
