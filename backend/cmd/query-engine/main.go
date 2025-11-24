package main

import (
	"context"
	"log"
	"search-engine/pkg/database"
	"search-engine/pkg/models"
	"search-engine/pkg/queue"
	"time"
)

func main() {
	// Initialize DB with authentication
	db, err := database.Connect("mongodb://rootuser:rootpassword@localhost:27017", "searchengine")
	if err != nil {
		log.Fatal(err)
	}

	// Example: Inserting a crawled page (Spider Service)
	newPage := models.Page{
		URL:       "https://example.com",
		Title:     "Example Domain",
		Content:   "This is the content...",
		CrawledAt: time.Now(),
	}

	// Use standard driver methods
	_, err = db.Database.Collection("pages").InsertOne(context.TODO(), newPage)

	// Initialize Redis
	q, err := queue.NewQueueClient("localhost:6379", "", 0)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// Seed the spider with diverse websites (The first seeds)
	// If we don't got this, the spider has nothing to start with
	seedURLs := []string{
		"https://news.ycombinator.com/", // Tech news
		"https://www.reddit.com/",       // Social media
		"https://www.bbc.com/",          // News
		"https://stackoverflow.com/",    // Tech Q&A
		"https://github.com/",           // Code repositories
		"https://www.wikipedia.org/",    // Encyclopedia
	}

	seededCount := 0
	for _, seedURL := range seedURLs {
		added, err := q.PushToSpider(ctx, seedURL)
		if err != nil {
			log.Printf("Error seeding %s: %v", seedURL, err)
		} else if added {
			log.Printf("Seeded spider with: %s", seedURL)
			seededCount++
		} else {
			log.Printf("URL already visited: %s", seedURL)
		}
	}
	log.Printf("Successfully seeded %d URLs", seededCount)

	// Check the Queue depth
	count, _ := q.QueueLength(ctx, queue.QueueSpider)
	log.Printf("Spider Queue Depth: %d", count)

}
