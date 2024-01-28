package models

import "gorm.io/gorm"

const (
	THREE_THREE_SEAT_LAYOUT = "3-3"
)

type Airplane struct {
	gorm.Model
	AirplaneModel string `json:"airplaneModel"`
	Manufacturer  string `json:"manufacturer"`
	Capacity      int    `json:"capacity"`
	SeatConfig    string `json:"seatConfig"`
	Entertainment bool   `json:"entertainment"`
	WiFi          bool   `json:"wifi"`
	PowerOutlets  bool   `json:"powerOutlets"`

	AirlineID uint `json:"-"`

	Airline         Airline          `json:"airline"`
	FlightSchedules []FlightSchedule `gorm:"foreignKey:AirplaneID" json:"-"`
}
