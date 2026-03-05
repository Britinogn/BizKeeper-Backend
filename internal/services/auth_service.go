package services

import (
	"context"
	"errors"
	"strings"

	"github.com/britinogn/bizkeeper/internal/model"
	"github.com/britinogn/bizkeeper/pkg/utils"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidInput           = errors.New("invalid registration data")
	ErrMissingRequiredFields  = errors.New("name, email, and password are required")
	ErrMissingLoginFields     = errors.New("email and password are required")
	ErrWeakPassword           = errors.New("password is too weak")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInvalidToken           = errors.New("invalid token")
	ErrDatabaseOperation      = errors.New("database operation failed")
	ErrUserNotFound           = errors.New("user not found")
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type AuthService struct {
	userRepo UserRepo
}

func NewAuthService(userRepo UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegistrationRequest) (*model.User, error) {
	if req == nil {
		return nil, ErrInvalidInput
	}

	// normalize users
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// validate user
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return nil, ErrMissingRequiredFields
	}

	// Optional: very basic extra rules for password strength
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	if !strings.Contains(req.Email, "@") {
		return nil, errors.New("invalid email format")
	}

	// check if email is already registered
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyRegistered
	}

	// create user
	user := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	if req == nil {
		return nil, "", ErrInvalidCredentials
	}

	// normalize users
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// validate user
	if req.Email == "" || req.Password == "" {
		return nil, "", ErrMissingLoginFields
	}

	// check if user exists
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// verify password
	match, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil || !match {
		return nil, "", ErrInvalidCredentials
	}
	user.Password = ""

	// generate token
	token, err := utils.GenerateToken(user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		return nil, "", ErrDatabaseOperation
	}
	return user, token, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	if req == nil {
		return nil, ErrInvalidInput
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrDatabaseOperation
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		if len(req.Password) < 8 {
			return nil, ErrWeakPassword
		}
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, ErrDatabaseOperation
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, ErrDatabaseOperation
	}

	return user, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrDatabaseOperation
	}

	return s.userRepo.DeleteUser(ctx, user)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}


// func (s *AuthService) GetProfile(ctx context.Context, email string) (*model.User, error) {
// 	user, err := s.userRepo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return nil, ErrUserNotFound
// 	}

// 	user.Password = ""

// 	return user, nil
// }
