package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleOwner Role = "owner"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"      json:"id"`
	FirstName string         `gorm:"not null"                  json:"first_name"`
	LastName  string         `gorm:"not null"                  json:"last_name"`
	Email     string         `gorm:"uniqueIndex;not null"      json:"email"`
	Password  string         `gorm:"not null"                  json:"-"`
	Role      Role           `gorm:"type:varchar(10);default:'owner'" json:"role"`
	CreatedAt time.Time      `                                 json:"created_at"`
	UpdatedAt time.Time      `                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                     json:"-"`
}

type RegistrationRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name"  validate:"required,min=2,max=50"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8"`
}


func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}