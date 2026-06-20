package models

import (
	"time"

	"github.com/google/uuid"
)

type Harvest struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PlantingID  uuid.UUID  `gorm:"index"`
	HarvestDate time.Time  `gorm:"not null"`
	Quantity    float64    `gorm:"not null"`
	Unit        string     `gorm:"not null"`
	QualityGrade string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationship
	Planting Planting `gorm:"foreignKey:PlantingID"`
}