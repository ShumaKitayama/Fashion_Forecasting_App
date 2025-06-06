package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/trendscout/backend/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInternalError      = errors.New("internal server error")
)

// TokenType defines the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	UserID  int    `json:"user_id"`
	TokenID string `json:"token_id,omitempty"`
	jwt.RegisteredClaims
}

// Service provides authentication related operations
type Service struct{}

// NewService creates a new authentication service
func NewService() *Service {
	return &Service{}
}

// GenerateTokens generates a new access and refresh token pair
func (s *Service) GenerateTokens(ctx context.Context, user *models.User) (accessToken string, refreshToken string, err error) {
	// Generate access token
	accessToken, _, err = s.generateToken(user.ID, AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token and ID
	refreshToken, tokenID, err := s.generateToken(user.ID, RefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis keyed by token ID
	refreshExpiration := 7 * 24 * time.Hour // 7 days
	redisKey := fmt.Sprintf("auth:refresh:%s", tokenID)

	err = models.SetWithTTL(ctx, redisKey, user.ID, refreshExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, err error) {
	// Get user by email
	user, err := models.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", "", ErrInvalidCredentials
	}

	// Verify password
	if !user.Authenticate(password) {
		return "", "", ErrInvalidCredentials
	}

	// Generate tokens
	return s.GenerateTokens(ctx, user)
}

// RefreshTokens refreshes an access token using a refresh token
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken string, err error) {
	// Verify refresh token
	claims, err := s.verifyToken(refreshToken, RefreshToken)
	if err != nil {
		return "", err
	}

	// Ensure refresh token exists in Redis (not logged out)
	redisKey := fmt.Sprintf("auth:refresh:%s", claims.TokenID)
	exists, err := models.KeyExists(ctx, redisKey)
	if err != nil {
		return "", ErrInternalError
	}
	if !exists {
		return "", ErrInvalidToken
	}

	// Get user
	user, err := models.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	// Generate new access token
	accessToken, _, err := s.generateToken(user.ID, AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// Logout invalidates a refresh token
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// Verify refresh token
	claims, err := s.verifyToken(refreshToken, RefreshToken)
	if err != nil {
		return err
	}

	// Delete refresh token from Redis
	redisKey := fmt.Sprintf("auth:refresh:%s", claims.TokenID)
	return models.Del(ctx, redisKey)
}

// VerifyAccessToken verifies an access token and returns the user ID
func (s *Service) VerifyAccessToken(accessToken string) (int, error) {
	claims, err := s.verifyToken(accessToken, AccessToken)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// generateToken generates a JWT token
func (s *Service) generateToken(userID int, tokenType TokenType) (string, string, error) {
	var expirationTime time.Time
	var tokenID string

	// Set expiration based on token type
	if tokenType == AccessToken {
		expirationTime = time.Now().Add(15 * time.Minute)
	} else {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
		tokenID = uuid.NewString()
	}

	// Create claims
	claims := &Claims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", "", errors.New("JWT_SECRET not set")
	}

	// Sign and get the complete encoded token as a string
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return signed, tokenID, nil
}

// verifyToken verifies a JWT token and returns the claims
func (s *Service) verifyToken(tokenString string, tokenType TokenType) (*Claims, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Basic type validation
	if tokenType == RefreshToken && claims.TokenID == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
