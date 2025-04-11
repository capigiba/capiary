package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// extractToken extracts the token from the Authorization header or query parameter.
func extractToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if len(token) == 0 {
		token = ctx.Query("Authorization")
	}
	return strings.TrimPrefix(token, "Bearer ")
}

// GenerateToken creates a JWT token for a user.
func (am *AuthUserMiddleware) GenerateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(am.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserByToken extracts user information from a JWT token.
func (am *AuthUserMiddleware) GetUserByToken(tokenStr string) (*entity.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(am.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Safely get userID
	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID type in token")
	}

	userID := uint64(userIDFloat)
	context := context.Background()
	user, err := am.userRepo.GetUserByID(context, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
