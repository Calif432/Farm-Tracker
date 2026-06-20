package models

import (
	"time"

	"github.com/google/uuid"
)

type Death struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AnimalID       uuid.UUID  `gorm:"uniqueIndex;not null"`
	DeathDate      time.Time  `gorm:"not null"`
	Reason         string     `gorm:"not null"`
	PreventionNotes string
	RecordedBy     *uuid.UUID `gorm:"index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relationships
	Animal   Animal `gorm:"foreignKey:AnimalID"`
	RecordedByUser *User `gorm:"foreignKey:RecordedBy"`
}