package services

import (
	"context"
	"errors"
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/middleware"
	"github.com/capigiba/capiary/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email, password string) (string, uint64, constant.Role, error)
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
	UpdateUser(ctx context.Context, user *entity.User) error
	UpdateAvatar(ctx context.Context, userID uint64, avatarPath, avatarFolder string) error
	GetUserByID(ctx context.Context, userID uint64) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
}

type userService struct {
	repo repositories.UserRepository
	auth middleware.MiddlewareInterface
}

// NewUserService returns a new user service.
func NewUserService(repo repositories.UserRepository, auth middleware.MiddlewareInterface) UserService {
	return &userService{
		repo: repo,
		auth: auth,
	}
}

// RegisterUser creates a new user with hashed password
func (s *userService) RegisterUser(ctx context.Context, user *entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if user.Email == "" || user.FirstName == "" || user.LastName == "" {
		return errors.New("missing necessary fields")
	}
	user.Password = string(hashedPassword)
	user.Status = constant.StatusPending // require accept from admin to use system
	user.Role = constant.RoleBasic
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.CreateUser(ctx, user)
}

// Login handles user login
func (s *userService) Login(ctx context.Context, email, password string) (string, uint64, constant.Role, error) {
	token, userID, role, err := s.auth.Login(email, password)
	return token, userID, role, err
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Compare old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password does not match")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdateUserPassword(ctx, userID, string(hashedPassword))
}

// UpdateUser updates user information (except avatar)
func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return s.repo.UpdateUser(ctx, user)
}

// UpdateAvatar updates the user's avatar
func (s *userService) UpdateAvatar(ctx context.Context, userID uint64, avatarPath, avatarFolder string) error {
	return s.repo.UpdateUserAvatar(ctx, userID, avatarPath, avatarFolder)
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, userID uint64) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// DeleteUser calls the repository to delete the user (soft or hard delete).
func (s *userService) DeleteUser(ctx context.Context, userID uint64) error {
	// If you want to do a soft-delete, you'd do:
	return s.repo.SoftDeleteUser(ctx, userID)
	// If you want to do a hard delete:
	// return s.repo.DeleteUser(ctx, userID)
}
