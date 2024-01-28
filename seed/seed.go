package seed

import (
	"fmt"

	"github.com/Laevateinn17/travelohi-backend/controllers"
	"github.com/Laevateinn17/travelohi-backend/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	if err := seedAirport(db); err != nil {
		fmt.Println(err.Error())
	}

	if err := seedAirline(db); err != nil {
		fmt.Println(err.Error())
	}

	if err := seedAirplane(db); err != nil {
		fmt.Println(err.Error())
	}
}

func seedAirport(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&models.Airport{}); err != nil {
		fmt.Printf("Error dropping Airport table: %v\n", err)
		return err
	}

	if err := db.AutoMigrate(&models.Airport{}); err != nil {
		fmt.Printf("Error creating Airport table: %v\n", err)
		return err
	}
	airports := []models.Airport{
		{
			Name:    "Mohamed Boudiaf International Airport",
			Code:    "CZL",
			City:    "Constantine",
			Country: "Algeria",
		},
		{
			Name:    "Chlef International Airport",
			Code:    "CFK",
			City:    "Chlef",
			Country: "Algeria",
		},
		{
			Name:    "Halim Perdanakusuma International Airport",
			Code:    "HLP",
			City:    "Jakarta",
			Country: "Indonesia",
		},
		{
			Name:    "Soekarno–Hatta International Airport",
			Code:    "CGK",
			City:    "Jakarta",
			Country: "Indonesia",
		},
		{
			Name:    "Hang Nadim International Airport",
			Code:    "BTH",
			City:    "Batam",
			Country: "Indonesia",
		},
		{
			Name:    "Kualanamu International Airport",
			Code:    "KNO",
			City:    "Medan",
			Country: "Indonesia",
		},
		{
			Name:    "Sultan Syarif Kasim II International Airport",
			Code:    "PKU",
			City:    "Pekanbaru",
			Country: "Indonesia",
		},
		{
			Name:    "Husein Sastranegara International Airport",
			Code:    "BDO",
			City:    "Bandung",
			Country: "Indonesia",
		},
	}

	for _, airport := range airports {
		if err := db.Create(&airport).Error; err != nil {
			fmt.Println("error: ", err)
			return err
			// return fmt.Errorf("error seeding airport %s", airport.Name)
		}
	}

	return nil
}

func seedAirline(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&models.Airline{}); err != nil {
		fmt.Printf("Error dropping Airport table: %v\n", err)
		return err
	}

	if err := db.AutoMigrate(&models.Airline{}); err != nil {
		fmt.Printf("Error creating Airport table: %v\n", err)
		return err
	}

	airlines := []models.Airline{
		{
			Name: "AirAsia Indonesia",
		},
		{
			Name: "Garuda Indonesia",
		},
		{
			Name: "Citilink",
		},
		{
			Name: "Batik Air",
		},
		{
			Name: "Super Air Jet",
		},
		{
			Name: "Lion Air",
		},
	}

	for _, airline := range airlines {
		if err := db.Create(&airline).Error; err != nil {
			fmt.Println("error: ", err)
			return err
		}
	}

	return nil

}

func seedAirplane(db *gorm.DB) error {
	airlines, err := controllers.GetAirlines(db)

	if err := db.Migrator().DropTable(&models.Airplane{}); err != nil {
		fmt.Printf("Error dropping Airport table: %v\n", err)
		return err
	}

	if err := db.AutoMigrate(&models.Airplane{}); err != nil {
		fmt.Printf("Error creating Airport table: %v\n", err)
		return err
	}

	if err != nil {
		return fmt.Errorf("error retrieving airlines")
	}

	airplanes := []models.Airplane{
		{
			AirplaneModel: "A320",
			Manufacturer:  "Airbus",
			Capacity:      168,
			SeatConfig:    models.THREE_THREE_SEAT_LAYOUT,
			Entertainment: true,
			WiFi:          true,
			PowerOutlets:  true,
		},
		{
			AirplaneModel: "737",
			Manufacturer:  "Boeing",
			Capacity:      162,
			SeatConfig:    models.THREE_THREE_SEAT_LAYOUT,
			Entertainment: false,
			WiFi:          false,
			PowerOutlets:  true,
		},
	}

	for _, airline := range airlines {
		for _, airplane := range airplanes {
			airplane.Airline = *airline
			if err := db.Create(&airplane).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
