package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the database
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Do not include in JSON responses
	CreatedAt    time.Time `json:"created_at"`
}

// CreateUser creates a new user in the database
func CreateUser(ctx context.Context, email, password string) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create the user in the database
	var user User
	err = PgPool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at`,
		email, hashedPassword).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := PgPool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No user found with this email
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := PgPool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No user found with this ID
		}
		return nil, err
	}

	return &user, nil
}

// Authenticate verifies a password against the user's hashed password
func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdatePassword updates a user's password
func (u *User) UpdatePassword(ctx context.Context, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password in the database
	_, err = PgPool.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		hashedPassword, u.ID)
	if err != nil {
		return err
	}

	// Update the local object
	u.PasswordHash = string(hashedPassword)
	return nil
} 