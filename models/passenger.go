package models

type Passenger struct {
	ID          string `json:"id"`
	TicketID    string `json:"ticketID" gorm:"foreignKey:ID"`
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"dateOfBirth"`
}
