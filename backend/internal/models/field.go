package models

import (
	"time"

	"github.com/google/uuid"
)

type Field struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name         string    `gorm:"not null"`
	SizeHectares float64
	SoilType     string
	Location     string
	FarmID       uuid.UUID `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relationships
	Farm       Farm       `gorm:"foreignKey:FarmID"`
	Plantings  []Planting `gorm:"foreignKey:FieldID"`
}