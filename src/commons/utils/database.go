package commons

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnect() {
	db, err := gorm.Open(postgres.Open(os.Getenv("GOOSE_DBSTRING")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can't connect to database")
	}
	DB = db
}
