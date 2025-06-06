package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/views"
)

// AuthController handles auth-related requests
type AuthController struct {
	authService *auth.Service
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *auth.Service) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register handles user registration
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existingUser, err := models.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	if existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Create new user
	user, err := models.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return user data (excluding password)
	ctx.JSON(http.StatusCreated, views.NewUserResponse(user))
}

// Login handles user login
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	accessToken, refreshToken, err := c.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		}
		return
	}

	// Return tokens
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh handles token refresh
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Refresh access token
	accessToken, err := c.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		switch err {
		case auth.ErrInvalidToken:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		case auth.ErrExpiredToken:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		}
		return
	}

	// Return new access token
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// Logout handles user logout
func (c *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	// Extract token
	tokenParts := authHeader[7:] // Remove "Bearer " prefix
	
	// Invalidate token
	if err := c.authService.Logout(ctx, tokenParts); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	ctx.Status(http.StatusOK)
} 