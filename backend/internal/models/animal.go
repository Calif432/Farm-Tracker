package models

import (
	"time"

	"github.com/google/uuid"
)

type Animal struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TagID       string     `gorm:"uniqueIndex;not null"`
	Type        string     `gorm:"not null;index"` // cow, sheep, goat, rabbit, chicken, etc.
	Breed       string
	Gender      string     `gorm:"check:gender IN ('male','female')"`
	DateOfBirth time.Time  `gorm:"not null"`
	Status      string     `gorm:"default:alive;check:status IN ('alive','dead','sold','slaughtered')"`
	ParentID    *uuid.UUID `gorm:"index"`
	MotherID    *uuid.UUID `gorm:"index"`
	FatherID    *uuid.UUID `gorm:"index"`
	FarmID      *uuid.UUID `gorm:"index"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	Parent     *Animal        `gorm:"foreignKey:ParentID"`
	Mother     *Animal        `gorm:"foreignKey:MotherID"`
	Father     *Animal        `gorm:"foreignKey:FatherID"`
	Farm       *Farm          `gorm:"foreignKey:FarmID"`
	Produce    []AnimalProduce `gorm:"foreignKey:AnimalID"`
	Vaccinations []Vaccination `gorm:"foreignKey:AnimalID"`
	Death      *Death         `gorm:"foreignKey:AnimalID"`
}