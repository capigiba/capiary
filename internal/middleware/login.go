package middleware

import (
	"errors"

	"github.com/capigiba/capiary/internal/domain/constant"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates the user and returns a JWT token.
func (am *AuthUserMiddleware) Login(email, password string) (string, uint64, constant.Role, error) {
	user, err := am.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", 0, "", errors.New("user not found")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", 0, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := am.GenerateToken(user)
	if err != nil {
		return "", 0, "", errors.New("failed to generate token")
	}

	return token, user.ID, user.Role, nil
}
