package entity

import (
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
)

type User struct {
	ID            uint64                 `json:"id" db:"id"`
	FirstName     string                 `json:"first_name" db:"first_name"`
	LastName      string                 `json:"last_name" db:"last_name"`
	UserName      string                 `json:"username" db:"username"`
	Email         string                 `json:"email" db:"email"`
	Password      string                 `json:"password" db:"password"`
	Status        constant.AccountStatus `json:"status" db:"status"`
	Role          constant.Role          `json:"role" db:"role"`
	Avatar        string                 `json:"avatar" db:"avatar"`
	AvatarFolder  string                 `json:"avatar_folder" db:"avatar_folder"`
	WalletBalance int64                  `json:"wallet_balance" db:"wallet_balance"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}
