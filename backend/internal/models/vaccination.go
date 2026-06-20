package models

import (
	"time"

	"github.com/google/uuid"
)

type Vaccination struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AnimalID        uuid.UUID  `gorm:"index"`
	VaccineName     string     `gorm:"not null"`
	AdministeredDate time.Time `gorm:"not null"`
	NextDueDate     *time.Time
	AdministeredBy  string
	BatchNumber     string
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relationship
	Animal Animal `gorm:"foreignKey:AnimalID"`
}