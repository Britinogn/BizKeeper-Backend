package services

import (
	"context"
	"errors"
	"strings"

	"github.com/britinogn/bizkeeper/internal/model"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidInput           = errors.New("invalid registration data")
	ErrMissingRequiredFields = errors.New("name, email, and password are required")
	ErrWeakPassword           = errors.New("password is too weak")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInvalidToken           = errors.New("invalid token")
	ErrDatabaseOperation      = errors.New("database operation failed")
)


type UserRepo interface{
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
}

type AuthService struct {
	userRepo UserRepo
}

func NewAuthService(userRepo UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}


func (s *AuthService) Register(ctx context.Context, req *model.RegistrationRequest) error {
	if req == nil {
		return ErrInvalidInput
	}

	// normalize users 
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// validate user
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return ErrMissingRequiredFields
	}

	// Optional: very basic extra rules for password strength
	if len(req.Password) < 8 {
		return ErrWeakPassword
	}

	if !strings.Contains(req.Email, "@") {
		return errors.New("invalid email format")
	}
	
	// check if email is already registered
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return ErrEmailAlreadyRegistered
	}

	// create user
	user := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	return s.userRepo.CreateUser(ctx, user)
}


func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) error {
	return nil
}