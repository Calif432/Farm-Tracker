package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Type          string     `gorm:"not null;check:type IN ('income','expense')"`
	Amount        float64    `gorm:"not null"`
	Category      string     `gorm:"not null"` // animal_sale, crop_sale, feed_purchase, labor, etc.
	Description   string
	TransactionDate time.Time `gorm:"not null"`
	RelatedToType string     // animal, planting, inventory, etc.
	RelatedToID   *uuid.UUID `gorm:"index"`
	RecordedBy    *uuid.UUID `gorm:"index"`
	ReceiptURL    string     // for storing image of receipt
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	RecordedByUser *User `gorm:"foreignKey:RecordedBy"`
}