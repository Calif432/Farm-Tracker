package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Password  string     `gorm:"not null"`
	FullName  string     `gorm:"not null"`
	Role      string     `gorm:"default:worker;check:role IN ('owner','manager','worker')"`
	FarmID    *uuid.UUID `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Farm *Farm `gorm:"foreignKey:FarmID"`
}