package middleware

import (
	"net/http"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/gin-gonic/gin"
)

// Auth is a middleware function that authenticates the user if a token is present.
func (am *AuthUserMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if len(token) == 0 {
			ctx.Next()
			return
		}

		userInfo, err := am.GetUserByToken(token)
		if err != nil || userInfo == nil {
			ctx.Next()
			return
		}

		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

// MustAuth ensures the user is authenticated; otherwise, returns an error.
func (am *AuthUserMiddleware) MustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if len(token) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userInfo, err := am.GetUserByToken(token)
		if err != nil || userInfo == nil || userInfo.Status == constant.StatusInactive {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}
