package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trendscout/backend/internal/auth"
	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/scraper"
)

// DataController handles data collection requests
type DataController struct{}

// NewDataController creates a new data controller
func NewDataController() *DataController {
	return &DataController{}
}

// CollectKeywordData handles data collection for a specific keyword
func (c *DataController) CollectKeywordData(ctx *gin.Context) {
	// Get keyword ID from path
	keywordIDStr := ctx.Param("id")
	keywordID, err := strconv.Atoi(keywordIDStr)
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

	// Check if keyword belongs to user
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keyword"})
		return
	}

	if keyword == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}

	if keyword.UserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Create timeout context for scraping
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Minute)
	defer cancel()

	// Initialize scraper service
	scraperService := scraper.NewService()

	// Perform scraping
	items, err := scraperService.ScrapeKeyword(timeoutCtx, keyword.Keyword)
	fmt.Printf("Scraping completed for keyword '%s': %d items found\n", keyword.Keyword, len(items))
	if err != nil {
		fmt.Printf("Scraping error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to collect data",
			"details": err.Error(),
		})
		return
	}

	// Convert scraped items to trend records and save to database
	totalVolume := len(items)
	if totalVolume > 0 {
		// Calculate sentiment (simple average for now)
		var totalSentiment float64
		for _, item := range items {
			// Simple sentiment calculation based on keywords
			sentiment := calculateSimpleSentiment(item.Title + " " + item.Content)
			totalSentiment += sentiment
		}
		avgSentiment := totalSentiment / float64(totalVolume)

		// Create trend record for today
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		
		// Save trend record using the correct function
		fmt.Printf("Saving trend record: KeywordID=%d, Volume=%d, Sentiment=%.2f\n", keywordID, totalVolume, avgSentiment)
		_, err := models.CreateTrendRecord(timeoutCtx, keywordID, today, totalVolume, avgSentiment)
		if err != nil {
			fmt.Printf("Failed to save trend record: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save trend data",
				"details": err.Error(),
			})
			return
		}
		fmt.Printf("Trend record saved successfully\n")

		// Save images to MongoDB
		imageCount := 0
		for _, item := range items {
			if item.ImageURL != "" {
				imageRecord := &models.Image{
					KeywordID: keywordID,
					ImageURL:  item.ImageURL,
					Caption:   item.Title + " " + item.Content,
					Tags:      item.Tags,
					FetchedAt: now,
				}
				
				if err := models.CreateImage(timeoutCtx, imageRecord); err != nil {
					fmt.Printf("Failed to save image %s: %v\n", item.ImageURL, err)
					// Log error but continue processing
					// ctx.JSON() でエラーを返すとループが止まる
					continue
				}
				imageCount++
			}
		}
		fmt.Printf("Saved %d images to MongoDB\n", imageCount)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data collection completed",
		"keyword": keyword.Keyword,
		"items_collected": totalVolume,
		"date": time.Now().Format("2006-01-02"),
	})
}

// calculateSimpleSentiment performs a simple sentiment analysis
func calculateSimpleSentiment(text string) float64 {
	// Simple sentiment analysis based on keywords
	positiveWords := []string{"good", "great", "excellent", "amazing", "love", "beautiful", "perfect", "best", "awesome", "fantastic"}
	negativeWords := []string{"bad", "terrible", "awful", "hate", "ugly", "worst", "horrible", "disgusting", "disappointing"}
	
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range positiveWords {
		if contains(text, word) {
			positiveCount++
		}
	}
	
	for _, word := range negativeWords {
		if contains(text, word) {
			negativeCount++
		}
	}
	
	if positiveCount > negativeCount {
		return 0.8 // Positive sentiment
	} else if negativeCount > positiveCount {
		return 0.3 // Negative sentiment
	}
	
	return 0.5 // Neutral sentiment
}

// contains checks if a string contains a substring (case insensitive)
func contains(text, substr string) bool {
	// Simple case-insensitive check
	return len(text) >= len(substr) && (text == substr || 
		(len(text) > len(substr) && containsHelper(text, substr)))
}

func containsHelper(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
} 