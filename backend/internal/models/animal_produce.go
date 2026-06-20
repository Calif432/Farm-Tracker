package models

import (
	"time"

	"github.com/google/uuid"
)

type AnimalProduce struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AnimalID    uuid.UUID  `gorm:"index"`
	ProduceType string     `gorm:"not null"` // milk, eggs, wool, meat, etc.
	Quantity    float64    `gorm:"not null"`
	Unit        string     `gorm:"not null"` // liters, kg, pieces, etc.
	RecordedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationship
	Animal Animal `gorm:"foreignKey:AnimalID"`
}