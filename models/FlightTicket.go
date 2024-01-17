package models


type FlightTicket struct {
	ID string `json:"id"`
	FlightID string `json:"flightID" gorm:"foreignKey:ID"`
	Passengers []Passenger `json:"passengers" gorm:"foreignKey:TicketID"`
}