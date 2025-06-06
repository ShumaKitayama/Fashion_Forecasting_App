package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Keyword represents a keyword in the database
type Keyword struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateKeyword adds a new keyword for a user
func CreateKeyword(ctx context.Context, userID int, keyword string) (*Keyword, error) {
	var k Keyword
	err := PgPool.QueryRow(ctx,
		`INSERT INTO keywords (user_id, keyword) VALUES ($1, $2) RETURNING id, user_id, keyword, created_at`,
		userID, keyword).Scan(&k.ID, &k.UserID, &k.Keyword, &k.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

// GetKeywordsForUser retrieves all keywords for a specific user
func GetKeywordsForUser(ctx context.Context, userID int) ([]*Keyword, error) {
	rows, err := PgPool.Query(ctx,
		`SELECT id, user_id, keyword, created_at FROM keywords WHERE user_id = $1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*Keyword
	for rows.Next() {
		var k Keyword
		if err := rows.Scan(&k.ID, &k.UserID, &k.Keyword, &k.CreatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, &k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}

// GetKeywordByID retrieves a keyword by its ID
func GetKeywordByID(ctx context.Context, id int) (*Keyword, error) {
	var k Keyword
	err := PgPool.QueryRow(ctx,
		`SELECT id, user_id, keyword, created_at FROM keywords WHERE id = $1`,
		id).Scan(&k.ID, &k.UserID, &k.Keyword, &k.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No keyword found with this ID
		}
		return nil, err
	}

	return &k, nil
}

// GetKeywordByName retrieves a keyword by its name (keyword text)
func GetKeywordByName(ctx context.Context, keyword string) (*Keyword, error) {
	var k Keyword
	err := PgPool.QueryRow(ctx,
		`SELECT id, user_id, keyword, created_at FROM keywords WHERE keyword = $1`,
		keyword).Scan(&k.ID, &k.UserID, &k.Keyword, &k.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No keyword found with this name
		}
		return nil, err
	}

	return &k, nil
}

// UpdateKeyword updates a keyword
func UpdateKeyword(ctx context.Context, id int, newKeyword string) error {
	result, err := PgPool.Exec(ctx,
		`UPDATE keywords SET keyword = $1 WHERE id = $2`,
		newKeyword, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("keyword not found")
	}

	return nil
}

// DeleteKeyword deletes a keyword
func DeleteKeyword(ctx context.Context, id int) error {
	result, err := PgPool.Exec(ctx,
		`DELETE FROM keywords WHERE id = $1`,
		id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("keyword not found")
	}

	return nil
}

// GetAllKeywords retrieves all keywords from the database
func GetAllKeywords(ctx context.Context) ([]*Keyword, error) {
	rows, err := PgPool.Query(ctx,
		`SELECT id, user_id, keyword, created_at FROM keywords ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*Keyword
	for rows.Next() {
		var k Keyword
		if err := rows.Scan(&k.ID, &k.UserID, &k.Keyword, &k.CreatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, &k)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
} 