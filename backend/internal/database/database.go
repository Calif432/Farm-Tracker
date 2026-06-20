package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Calif432/farmtrack/backend/internal/config"
	"github.com/Calif432/farmtrack/backend/internal/models"
)

var DB *gorm.DB

func Connect() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Database connected successfully!")

	// AutoMigrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Farm{},
		&models.Animal{},
		&models.AnimalProduce{},
		&models.Vaccination{},
		&models.Death{},
		&models.Field{},
		&models.Planting{},
		&models.Harvest{},
		&models.InventoryItem{},
		&models.Transaction{},
		&models.Task{},
	)
	if err != nil {
		log.Printf("⚠️ Auto migration warning: %v", err)
	} else {
		log.Println("✅ All models migrated successfully!")
	}
}