package models

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name          string     `gorm:"not null;index"`
	Category      string     `gorm:"not null;index"`
	Quantity      float64    `gorm:"default:0"`
	Unit          string
	MinStockLevel float64    `gorm:"default:0"`
	CostPerUnit   float64
	Supplier      string
	Location      string
	FarmID        *uuid.UUID `gorm:"index"`  // Changed to pointer (optional)
	LastUpdated   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationship
	Farm Farm `gorm:"foreignKey:FarmID"`
}