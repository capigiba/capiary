package middleware

import (
	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

// MiddlewareInterface is a placeholder interface – define your needed methods
type MiddlewareInterface interface {
	Auth() gin.HandlerFunc
	MustAuth() gin.HandlerFunc
	Login(email, password string) (string, uint64, constant.Role, error)
	GetUserByToken(tokenStr string) (*entity.User, error)
}

// UserRepository defines the methods needed from your user storage layer
type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID uint64) (*entity.User, error)
}

// AuthUserMiddleware handles user authentication
type AuthUserMiddleware struct {
	userRepo  UserRepository
	secretKey string
}

// NewAuthUserMiddleware creates a new AuthUserMiddleware
func NewAuthUserMiddleware(repo UserRepository, secretKey string) *AuthUserMiddleware {
	return &AuthUserMiddleware{
		userRepo:  repo,
		secretKey: secretKey,
	}
}
