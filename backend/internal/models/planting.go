package models

import (
	"time"

	"github.com/google/uuid"
)

type Planting struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FieldID            uuid.UUID  `gorm:"index"`
	CropType           string     `gorm:"not null;index"` // maize, beans, tomatoes, tea, coffee, etc.
	Variety            string
	PlantingDate       time.Time  `gorm:"not null"`
	ExpectedHarvestDate *time.Time
	ActualHarvestDate  *time.Time
	Status             string     `gorm:"default:growing;check:status IN ('planned','growing','harvested','failed')"`
	QuantityPlanted    float64
	Unit               string     // kg, bags, etc.
	Notes              string
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relationships
	Field    Field     `gorm:"foreignKey:FieldID"`
	Harvests []Harvest `gorm:"foreignKey:PlantingID"`
}