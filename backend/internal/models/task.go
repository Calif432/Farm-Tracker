package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title          string     `gorm:"not null"`
	Description    string
	DueDate        *time.Time `gorm:"index"`
	CompletedAt    *time.Time
	Status         string     `gorm:"default:pending;check:status IN ('pending','in_progress','completed','cancelled')"`
	Priority       string     `gorm:"default:medium;check:priority IN ('low','medium','high','urgent')"`
	AssignedTo     *uuid.UUID `gorm:"index"`
	RelatedToType  string     // animal, planting, task, etc.
	RelatedToID    *uuid.UUID
	CreatedBy      uuid.UUID  `gorm:"index"`
	ReminderDate   *time.Time
	IsRecurring    bool       `gorm:"default:false"`
	RecurrenceRule string     // daily, weekly, monthly
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relationships
	AssignedToUser *User `gorm:"foreignKey:AssignedTo"`
	CreatedByUser  User  `gorm:"foreignKey:CreatedBy"`
}