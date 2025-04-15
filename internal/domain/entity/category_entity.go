package entity

import (
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
)

type Category struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Access    []constant.Role `json:"access"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
