package views

import (
	"time"

	"github.com/trendscout/backend/internal/models"
)

// KeywordResponse represents the keyword data returned in API responses
type KeywordResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
}

// NewKeywordResponse creates a new keyword response from a keyword model
func NewKeywordResponse(keyword *models.Keyword) *KeywordResponse {
	return &KeywordResponse{
		ID:        keyword.ID,
		UserID:    keyword.UserID,
		Keyword:   keyword.Keyword,
		CreatedAt: keyword.CreatedAt,
	}
}

// KeywordListResponse represents a list of keywords
type KeywordListResponse struct {
	Keywords []*KeywordResponse `json:"keywords"`
	Count    int                `json:"count"`
}

// NewKeywordListResponse creates a new keyword list response
func NewKeywordListResponse(keywords []*models.Keyword) *KeywordListResponse {
	keywordResponses := make([]*KeywordResponse, len(keywords))
	for i, keyword := range keywords {
		keywordResponses[i] = NewKeywordResponse(keyword)
	}

	return &KeywordListResponse{
		Keywords: keywordResponses,
		Count:    len(keywordResponses),
	}
} 