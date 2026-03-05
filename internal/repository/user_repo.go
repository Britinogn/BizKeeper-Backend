package repository

import (
	"context"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/britinogn/bizkeeper/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	// hash password 
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// assign hashed password to user
	user.Password = hashedPassword

	// create user
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	return &user, r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// func (r *UserRepository) DeleteUser(ctx context.Context, user *model.User) error {
// 	return r.db.WithContext(ctx).Delete(user).Error
// }

func (r *UserRepository) DeleteUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Unscoped().Delete(user).Error
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	return &user, r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
}