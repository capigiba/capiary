package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uint64) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	SoftDeleteUser(ctx context.Context, userID uint64) error
	DeleteUser(ctx context.Context, userID uint64) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	UpdateUserPassword(ctx context.Context, userID uint64, hashedPassword string) error
	UpdateUserAvatar(ctx context.Context, userID uint64, avatarPath, avatarFolder string) error
}

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepo returns a concrete implementation of UserRepository backed by Postgres/sqlx.
func NewUserRepo(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

// CreateUser inserts a new user into the database.
func (r *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			first_name, last_name, username, email, password,
			status, role, avatar, avatar_folder, wallet_balance,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5,
		        $6, $7, $8, $9, $10,
		        $11, $12)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Email,
		user.Password,
		user.Status,
		user.Role,
		user.Avatar,
		user.AvatarFolder,
		user.WalletBalance,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	return err
}

// GetUserByEmail retrieves a user by their email address.
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, username, email,
			password, status, role, avatar, avatar_folder,
			wallet_balance, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *userRepo) GetUserByID(ctx context.Context, userID uint64) (*entity.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, username, email,
			password, status, role, avatar, avatar_folder,
			wallet_balance, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information (except password/avatar).
func (r *userRepo) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET 
			first_name     = $1,
			last_name      = $2,
			username       = $3,
			email          = $4,
			status         = $5,
			role           = $6,
			wallet_balance = $7,
			updated_at     = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.Email,
		user.Status,
		user.Role,
		user.WalletBalance,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

// SoftDeleteUser sets the user status to "deleted".
func (r *userRepo) SoftDeleteUser(ctx context.Context, userID uint64) error {
	query := `
		UPDATE users 
		SET status = $1 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, constant.StatusDeleted, userID)
	return err
}

// DeleteUser permanently removes the user from the database.
func (r *userRepo) DeleteUser(ctx context.Context, userID uint64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user with ID %d: %v", userID, err)
	}
	return nil
}

// UpdateUserPassword updates the hashed password for a user.
func (r *userRepo) UpdateUserPassword(ctx context.Context, userID uint64, hashedPassword string) error {
	query := `
		UPDATE users
		SET 
			password = $1,
			updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	return err
}

// UpdateUserAvatar updates the user's avatar and avatar folder.
func (r *userRepo) UpdateUserAvatar(ctx context.Context, userID uint64, avatarPath, avatarFolder string) error {
	query := `
		UPDATE users
		SET 
			avatar = $1,
			avatar_folder = $2,
			updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, avatarPath, avatarFolder, time.Now(), userID)
	return err
}

// GetAllUsers retrieves all users from the database.
func (r *userRepo) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, username, email,
			password, status, role, avatar, avatar_folder,
			wallet_balance, created_at, updated_at
		FROM users
	`
	var users []entity.User
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}
