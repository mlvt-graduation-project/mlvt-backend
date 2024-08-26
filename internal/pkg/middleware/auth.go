package middleware

import (
	"mlvt/internal/entity"
	"mlvt/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthUserMiddleware handles user authentication
type AuthUserMiddleware struct {
	authService *service.AuthService
}

// NewAuthUserMiddleware creates a new AuthUserMiddleware
func NewAuthUserMiddleware(authService *service.AuthService) *AuthUserMiddleware {
	return &AuthUserMiddleware{
		authService: authService,
	}
}

// Auth is a middleware function that authenticates the user if a token is present
func (am *AuthUserMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if len(token) == 0 {
			ctx.Next()
			return
		}

		userInfo, err := am.authService.GetUserByToken(token)
		if err != nil || userInfo == nil {
			ctx.Next()
			return
		}

		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

// MustAuth ensures the user is authenticated; otherwise, returns an error
func (am *AuthUserMiddleware) MustAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		if len(token) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userInfo, err := am.authService.GetUserByToken(token)
		if err != nil || userInfo == nil || userInfo.Status == entity.UserStatusSuspended || userInfo.Status == entity.UserStatusDeleted {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

// extractToken extracts the token from the Authorization header or query parameter
func extractToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if len(token) == 0 {
		token = ctx.Query("Authorization")
	}
	return strings.TrimPrefix(token, "Bearer ")
}
