package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Queue Names
const (
	QueueSpider  = "spider_queue"  // List: URLs to crawl
	QueueIndexer = "indexer_queue" // List: Document IDs to index
	SetVisited   = "visited_urls"  // Set: Deduplication
)

type QueueClient struct {
	Client *redis.Client
}

// NewQueueClient initializes the Redis connection
func NewQueueClient(addr string, password string, db int) (*QueueClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &QueueClient{Client: rdb}, nil
}

// --- Spider Operations ---

// PushToSpider adds a URL to the crawl queue only if it hasn't been visited
func (q *QueueClient) PushToSpider(ctx context.Context, url string) (bool, error) {
	// 1. Check if already visited
	isMember, err := q.Client.SIsMember(ctx, SetVisited, url).Result()
	if err != nil {
		return false, err
	}

	// If true, we have already seen this URL. Skip it.
	if isMember {
		return false, nil
	}

	// 2. Mark as visited
	if err := q.Client.SAdd(ctx, SetVisited, url).Err(); err != nil {
		return false, err
	}

	// 3. Push to Queue (Right Push)
	if err := q.Client.RPush(ctx, QueueSpider, url).Err(); err != nil {
		return false, err
	}

	return true, nil
}

// PopFromSpider gets the next URL to crawl (Blocking Pop)
func (q *QueueClient) PopFromSpider(ctx context.Context) (string, error) {
	// BLPop blocks until an item is available or timeout is reached (0 = infinite)
	result, err := q.Client.BLPop(ctx, 0, QueueSpider).Result()
	if err != nil {
		return "", err
	}
	// result[0] is the key name, result[1] is the value
	return result[1], nil
}

// --- Indexer Operations ---

// PushToIndex adds a Document ID (string) to the indexer queue
func (q *QueueClient) PushToIndex(ctx context.Context, docID string) error {
	return q.Client.RPush(ctx, QueueIndexer, docID).Err()
}

// PopFromIndex gets the next Document ID to process
func (q *QueueClient) PopFromIndex(ctx context.Context) (string, error) {
	result, err := q.Client.BLPop(ctx, 0, QueueIndexer).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil
}

// --- Utility ---

// QueueLength checks how many items are pending
func (q *QueueClient) QueueLength(ctx context.Context, queueName string) (int64, error) {
	return q.Client.LLen(ctx, queueName).Result()
}
