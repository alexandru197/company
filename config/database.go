package config

import (
	"fmt"
	"log"

	"github.com/alexandru197/company/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	// Migrate the schema.
	if err := db.AutoMigrate(&model.Company{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}
	return db, nil
}
