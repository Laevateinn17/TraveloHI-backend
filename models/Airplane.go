package models

import "gorm.io/gorm"

type Airplane struct {
	gorm.Model
	Manufacturer  string `json:"manufacturer"`
	Capacity      int    `json:"capacity"`
	SeatConfig    string `json:"seatConfig"`
	Entertainment bool   `json:"entertainment"`
	WiFi          bool   `json:"wifi"`
	PowerOutlets  bool   `json:"powerOutlets"`
	FoodService   bool   `json:"foodService"`

	AirlineID uint `json:"-"`

	// Airline         Airline          `json:"airline"`
	FlightSchedules []FlightSchedule `gorm:"foreignKey:AirplaneID" json:"-"`
}
