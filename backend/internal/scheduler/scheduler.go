package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/trendscout/backend/internal/models"
	"github.com/trendscout/backend/internal/scraper"
)

// Service handles scheduled operations
type Service struct {
	scraperService *scraper.Service
	ticker         *time.Ticker
	quit           chan struct{}
}

// NewService creates a new scheduler service
func NewService() *Service {
	return &Service{
		scraperService: scraper.NewService(),
		quit:           make(chan struct{}),
	}
}

// Start begins the scheduled data collection
func (s *Service) Start() {
	log.Println("Starting data collection scheduler...")
	
	// Run immediately on startup
	go s.collectAllKeywordsData()
	
	// Schedule to run every 24 hours
	s.ticker = time.NewTicker(24 * time.Hour)
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Println("Running scheduled data collection...")
				s.collectAllKeywordsData()
			case <-s.quit:
				log.Println("Stopping data collection scheduler...")
				return
			}
		}
	}()
}

// Stop terminates the scheduler
func (s *Service) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.quit)
}

// collectAllKeywordsData collects data for all keywords in the system
func (s *Service) collectAllKeywordsData() {
	ctx := context.Background()
	
	// Get all keywords from the database
	keywords, err := models.GetAllKeywords(ctx)
	if err != nil {
		log.Printf("Failed to get keywords for scheduled collection: %v", err)
		return
	}

	log.Printf("Starting scheduled data collection for %d keywords", len(keywords))

	// Collect data for each keyword
	for _, keyword := range keywords {
		log.Printf("Collecting data for keyword: %s (ID: %d)", keyword.Keyword, keyword.ID)
		
		// Create context with timeout for each keyword
		keywordCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		
		items, err := s.scraperService.ScrapeKeyword(keywordCtx, keyword.Keyword)
		if err != nil {
			log.Printf("Failed to collect data for keyword %s: %v", keyword.Keyword, err)
			cancel()
			continue
		}
		
		log.Printf("Collected %d items for keyword %s", len(items), keyword.Keyword)
		cancel()
		
		// Rate limiting between keywords
		time.Sleep(2 * time.Second)
	}
	
	log.Printf("Completed scheduled data collection")
}

// ForceCollectionForKeyword manually triggers data collection for a specific keyword
func (s *Service) ForceCollectionForKeyword(ctx context.Context, keywordID int) error {
	keyword, err := models.GetKeywordByID(ctx, keywordID)
	if err != nil {
		return err
	}
	
	if keyword == nil {
		return fmt.Errorf("keyword not found")
	}
	
	log.Printf("Manual data collection triggered for keyword: %s (ID: %d)", keyword.Keyword, keyword.ID)
	
	items, err := s.scraperService.ScrapeKeyword(ctx, keyword.Keyword)
	if err != nil {
		return err
	}
	
	log.Printf("Manual collection completed for %s: %d items", keyword.Keyword, len(items))
	return nil
} 