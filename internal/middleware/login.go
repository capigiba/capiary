package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/capigiba/capiary/internal/domain/constant"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates the user and returns a JWT token.
func (am *AuthUserMiddleware) Login(email, password string) (string, uint64, constant.Role, error) {
	context := context.Background()
	user, err := am.userRepo.GetUserByEmail(context, email)
	if err != nil {
		resultErr := fmt.Sprint("user not found: ", err)
		return "", 0, "", errors.New(resultErr)
	}

	if user == nil {
		return "", 0, "", errors.New("user not found")
	}

	// Check if user.Password is empty
	if user.Password == "" {
		return "", 0, "", errors.New("user password not found")
	}

	// Ensure password is not empty before comparison
	if password == "" {
		return "", 0, "", errors.New("password cannot be empty")
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
