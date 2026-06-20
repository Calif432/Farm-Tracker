package models

import (
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name         string    `gorm:"not null"`
	Location     string
	SizeHectares float64
	OwnerID      uuid.UUID `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relationships
	Owner     User      `gorm:"foreignKey:OwnerID"`
	Animals   []Animal  `gorm:"foreignKey:FarmID"`
	Fields    []Field   `gorm:"foreignKey:FarmID"`
	Inventory []InventoryItem `gorm:"foreignKey:FarmID"`
}