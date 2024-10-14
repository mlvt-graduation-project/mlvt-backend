package middleware

import (
	"net/http"
	"time"

	"mlvt/internal/entity"

	"github.com/gin-gonic/gin"
)

// MockAuthMiddleware is a mock implementation of AuthUserMiddleware
type MockAuthMiddleware struct {
	// Optional: Add fields to customize behavior
}

// NewMockAuthMiddleware creates a new instance of MockAuthMiddleware
func NewMockAuthMiddleware() *MockAuthMiddleware {
	return &MockAuthMiddleware{}
}

// MustAuthAuthenticated simulates an authenticated user
func (m *MockAuthMiddleware) MustAuthAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inject a mock user into the context
		userInfo := &entity.User{
			ID:           1,
			FirstName:    "Mock",
			LastName:     "User",
			UserName:     "mockuser",
			Email:        "mockuser@example.com",
			Password:     "hashedpassword",
			Status:       1, // UserStatusAvailable
			Premium:      true,
			Role:         "User",
			Avatar:       "avatar.jpg",
			AvatarFolder: "avatars/",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		c.Set("userInfo", userInfo)
		c.Next()
	}
}

// MustAuthUnauthenticated simulates an unauthenticated request
func (m *MockAuthMiddleware) MustAuthUnauthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}

// AuthAuthenticated simulates an optional authenticated user
func (m *MockAuthMiddleware) AuthAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inject a mock user into the context
		userInfo := &entity.User{
			ID:           1,
			FirstName:    "Mock",
			LastName:     "User",
			UserName:     "mockuser",
			Email:        "mockuser@example.com",
			Password:     "hashedpassword",
			Status:       1, // UserStatusAvailable
			Premium:      true,
			Role:         "User",
			Avatar:       "avatar.jpg",
			AvatarFolder: "avatars/",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		c.Set("userInfo", userInfo)
		c.Next()
	}
}

// AuthUnauthenticated simulates an optional unauthenticated request
func (m *MockAuthMiddleware) AuthUnauthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
