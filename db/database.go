package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	database   *gorm.DB
	once       sync.Once
	connString = "host=localhost port=5432 user=postgres password=admin dbname=travelohi sslmode=disable"
)

func Connect() (*gorm.DB, error) {
	once.Do(func() {
		var err error
		database, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			fmt.Println("error connecting to database: ", err.Error())
		}
	})

	return database, nil
}
