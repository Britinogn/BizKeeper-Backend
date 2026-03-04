package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"     json:"id"`
	SessionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"session_id"`
	Name      string         `gorm:"not null"                 json:"name"`
	Quantity  int            `gorm:"not null"                 json:"quantity"`
	UnitPrice float64        `gorm:"not null"                 json:"unit_price"`
	Category  string         `gorm:"not null"                 json:"category"`
	Notes     string         `                                json:"notes"`
	CreatedAt time.Time      `                                json:"created_at"`
	UpdatedAt time.Time      `                                json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                    json:"-"`
}

func (pi *ProductItem) BeforeCreate(tx *gorm.DB) error {
	pi.ID = uuid.New()
	return nil
}