package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var connString = "host=localhost port=5432 user=postgres password=admin dbname=travelohi sslmode=disable"

func Connect() (*gorm.DB, error) {
	fmt.Println("--- Connecting to database... ---")

	database, err := gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return database, nil
}

