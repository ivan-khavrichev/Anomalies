package psql

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"team/transmitter/internal/domain"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("can`t connect to database")
		return nil, err
	}

	if err := db.AutoMigrate(&domain.AnomalyMessage{}); err != nil {
		log.Println("can`t automigrate database")
		return nil, err
	}

	return db, nil
}
