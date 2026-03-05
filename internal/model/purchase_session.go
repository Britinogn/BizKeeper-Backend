package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentTransfer PaymentMethod = "transfer"
	PaymentCredit   PaymentMethod = "credit"
	PaymentOther    PaymentMethod = "other"
)

type PurchaseSession struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey"     json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	PurchaseDate     time.Time      `gorm:"not null"                 json:"purchase_date"`
	SupplierName     string         `gorm:"not null"                 json:"supplier_name"`
	PaymentMethod    PaymentMethod  `gorm:"type:varchar(20);not null" json:"payment_method"`
	InvoiceReference *string        `                                json:"invoice_reference,omitempty"`
	Notes            *string        `                                json:"notes,omitempty"`
	ProductItems     []ProductItem  `gorm:"foreignKey:SessionID"     json:"product_items,omitempty"`
	CreatedAt        time.Time      `                                json:"created_at"`
	UpdatedAt        time.Time      `                                json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                    json:"-"`
}

func (ps *PurchaseSession) BeforeCreate(tx *gorm.DB) error {
	ps.ID = uuid.New()
	return nil
}
