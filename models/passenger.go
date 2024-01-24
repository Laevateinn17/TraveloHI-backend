package models

import "gorm.io/gorm"

type Passenger struct {
	gorm.Model
	FlightScheduleID uint   `json:"-"`
	FlightTicketID   uint   `json:"-"`
	FirstName        string `json:"firstName"`
	MiddleName       string `json:"middleName"`
	LastName         string `json:"lastName"`
	Gender           string `json:"gender"`
	DateOfBirth      string `json:"dateOfBirth"`
}
