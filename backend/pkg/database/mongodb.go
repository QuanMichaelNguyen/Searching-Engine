package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect initializes the MongoDB connection
func Connect(uri, dbName string) (*DBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	// Set a pool size to handle concurrent spider/indexer threads
	clientOptions.SetMinPoolSize(10)
	clientOptions.SetMaxPoolSize(100)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully")

	db := client.Database(dbName)
	instance := &DBClient{
		Client:   client,
		Database: db,
	}

	// Initialize Indexes immediately upon connection
	if err := instance.createIndexes(); err != nil {
		return nil, err
	}

	return instance, nil
}

// 	INDEX IMPLEMENTATION

// createIndexes ensures all necessary MongoDB indexes exist
func (db *DBClient) createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Pages Collection Indexes
	pages := db.Database.Collection("pages")

	// Create unique index on URL
	// URL must be unique to prevent duplicate crawls
	_, err := pages.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create page url index: %w", err)
	}

	// 2. Inverted Index Collection Indexes
	index := db.Database.Collection("index")

	// Create unique index on term
	// Fast lookup by Term
	_, err = index.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "term", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create term index: %w", err)
	}

	// 3. PageRank Collection Indexes
	pagerank := db.Database.Collection("pagerank")

	// Create unique index on URL
	// Fast lookup by URL
	_, err = pagerank.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Create non-unique index on score
	// Index on Score for faster sorting/retrieval of top pages
	_, err = pagerank.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "score", Value: -1}}, // Descending order
	})

	log.Println("Database indexes verified/created")
	return nil
}

// url, term, url must be unique inside pages, index, pagerank
// score in pagerank don't have to be unique
