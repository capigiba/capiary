package entity

import (
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
)

type User struct {
	ID            uint64                 `json:"id"`
	FirstName     string                 `json:"first_name"`
	LastName      string                 `json:"last_name"`
	UserName      string                 `json:"username"`
	Email         string                 `json:"email"`
	Password      string                 `json:"password"`
	Status        constant.AccountStatus `json:"status"`
	Role          constant.Role          `json:"role"`
	Avatar        string                 `json:"avatar"`        // file name on cloud
	AvatarFolder  string                 `json:"avatar_folder"` // Folder that contain the avatar image on cloud
	WalletBalance int64                  `json:"wallet_balance"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
