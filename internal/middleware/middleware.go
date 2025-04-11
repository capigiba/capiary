package middleware

import (
	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/repositories"
	"github.com/gin-gonic/gin"
)

// MiddlewareInterface is a placeholder interface – define your needed methods
type MiddlewareInterface interface {
	Auth() gin.HandlerFunc
	MustAuth() gin.HandlerFunc
	Login(email, password string) (string, uint64, constant.Role, error)
	GetUserByToken(tokenStr string) (*entity.User, error)
}

// AuthUserMiddleware handles user authentication
type AuthUserMiddleware struct {
	userRepo  repositories.UserRepository
	secretKey string
}

// NewAuthUserMiddleware creates a new AuthUserMiddleware
func NewAuthUserMiddleware(repo repositories.UserRepository, secretKey string) *AuthUserMiddleware {
	return &AuthUserMiddleware{
		userRepo:  repo,
		secretKey: secretKey,
	}
}
