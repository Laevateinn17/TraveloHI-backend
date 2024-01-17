package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func Connect(connString string) (*gorm.DB, error) {
	fmt.Println("--- Connecting to database... ---")

	database, err := gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return database, nil
}

