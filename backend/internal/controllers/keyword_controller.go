package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/views"
)

// KeywordController handles keyword-related requests
type KeywordController struct{}

// NewKeywordController creates a new keyword controller
func NewKeywordController() *KeywordController {
	return &KeywordController{}
}

// CreateKeywordRequest represents the request body for keyword creation
type CreateKeywordRequest struct {
	Keyword string `json:"keyword" binding:"required,min=2,max=100"`
}

// UpdateKeywordRequest represents the request body for keyword update
type UpdateKeywordRequest struct {
	Keyword string `json:"keyword" binding:"required,min=2,max=100"`
}

// CreateKeyword handles keyword creation
func (c *KeywordController) CreateKeyword(ctx *gin.Context) {
	var req CreateKeywordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create keyword
	keyword, err := models.CreateKeyword(ctx, userID, req.Keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keyword"})
		return
	}

	// Return keyword data
	ctx.JSON(http.StatusCreated, views.NewKeywordResponse(keyword))
}

// GetKeywords handles retrieving all keywords for a user
func (c *KeywordController) GetKeywords(ctx *gin.Context) {
	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get keywords
	keywords, err := models.GetKeywordsForUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keywords"})
		return
	}

	// Return keywords
	ctx.JSON(http.StatusOK, views.NewKeywordListResponse(keywords))
}

// UpdateKeyword handles keyword update
func (c *KeywordController) UpdateKeyword(ctx *gin.Context) {
	var req UpdateKeywordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get keyword ID from path
	keywordID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get keyword to check ownership
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	// Check if user owns the keyword
	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Update keyword
	if err := models.UpdateKeyword(ctx, keywordID, req.Keyword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update keyword"})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteKeyword handles keyword deletion
func (c *KeywordController) DeleteKeyword(ctx *gin.Context) {
	// Get keyword ID from path
	keywordID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid keyword ID"})
		return
	}

	// Get authenticated user ID
	userID, exists := auth.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get keyword to check ownership
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	// Check if user owns the keyword
	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Delete keyword
	if err := models.DeleteKeyword(ctx, keywordID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete keyword"})
		return
	}

	ctx.Status(http.StatusOK)
} 