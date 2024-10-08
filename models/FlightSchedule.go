package models

import (
	"time"

	"gorm.io/gorm"
)

type FlightSchedule struct {
	gorm.Model

	FlightNumber         string    `json:"flightNumber"`
	DepartureTime        time.Time `json:"departureTime"`
	ArrivalTime          time.Time `json:"arrivalTime"`
	DepartureAirportID   uint      `json:"-"`
	DestinationAirportID uint      `json:"-"`
	AirplaneID           uint      `json:"-"`
	// FlightTicketID       uint      `json:"-"`
	Price uint64 `json:"price"`

	//services

	FoodService bool `json:"foodService"`

	DepartureAirport   Airport        `json:"departureAirport"`
	DestinationAirport Airport        `json:"destinationAirport"`
	Airplane           Airplane       `json:"airplane"`
	Passengers         []Passenger    `json:"passengers" gorm:"foreignKey:FlightScheduleID"`
	FlightTickets      []FlightTicket `gorm:"many2many:flight_schedules_tickets"`
}
