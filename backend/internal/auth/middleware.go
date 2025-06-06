package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthUserKey is the key used to store the authenticated user ID in the context
const AuthUserKey = "auth_user_id"

// Middleware handles JWT authentication
type Middleware struct {
	authService *Service
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(authService *Service) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

// Authenticate is a middleware that authenticates requests with JWT
func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Verify the token
		userID, err := m.authService.VerifyAccessToken(tokenString)
		if err != nil {
			var status int
			var message string

			switch err {
			case ErrExpiredToken:
				status = http.StatusUnauthorized
				message = "Token has expired"
			case ErrInvalidToken:
				status = http.StatusUnauthorized
				message = "Invalid token"
			default:
				status = http.StatusInternalServerError
				message = "Internal server error"
			}

			c.JSON(status, gin.H{"error": message})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set(AuthUserKey, userID)
		c.Next()
	}
}

// GetUserID gets the authenticated user ID from the context
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(AuthUserKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int)
	return id, ok
} 