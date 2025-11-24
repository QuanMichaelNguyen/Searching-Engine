package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Page: a crawled website with local storage support and link tracking
// Collection: "pages"
// bson: used to store in database
// json: structure for JSON response/request
type Page struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL           string             `bson:"url" json:"url"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Content       string             `bson:"content" json:"content"`                 // Clean text content
	Language      string             `bson:"language" json:"language"`               // Detected language (ISO 639-1)
	LanguageScore float64            `bson:"language_score" json:"language_score"`   // Confidence score for language detection
	LocalFilePath string             `bson:"local_file_path" json:"local_file_path"` // Path to locally stored HTML file
	FileSize      int64              `bson:"file_size" json:"file_size"`             // Size of stored file in bytes
	ContentHash   string             `bson:"content_hash" json:"content_hash"`       // SHA256 hash of content for deduplication
	ReferrerURL   string             `bson:"referrer_url" json:"referrer_url"`       // URL of the page that linked to this page
	OutboundLinks []string           `bson:"outbound_links" json:"outbound_links"`   // URLs this page links to
	InboundCount  int                `bson:"inbound_count" json:"inbound_count"`     // Number of pages linking to this page
	Depth         int                `bson:"depth" json:"depth"`                     // Crawl depth from seed URLs
	Domain        string             `bson:"domain" json:"domain"`                   // Domain of the URL for easier querying
	CrawledAt     time.Time          `bson:"crawled_at" json:"crawled_at"`
	LastModified  time.Time          `bson:"last_modified" json:"last_modified"` // From HTTP headers
	Metadata      map[string]string  `bson:"metadata" json:"metadata"`           // Extra headers, etc
}

// InvertedIndexEntry represents a single term and the list of documents which has that term
// Collection: "index"
type InvertedIndexEntry struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Term     string             `bson:"term"`
	Postings []Posting          `bson:"postings"` // List of docs containing the term
}

// Posting represents a specific document's relationship to a term.
type Posting struct {
	DocID     primitive.ObjectID `bson:"doc_id"`
	Positions []int              `bson:"positions"` // Essential for phrase search (e.g. "Code" at pos 5, "Search" at pos 6)
	TF        float64            `bson:"tf"`        // Term Frequency score for this doc
}

// Link represents a directed link between two pages
// Collection: "links" - optimized for PageRank calculations
type Link struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FromURL    string             `bson:"from_url"`    // Source page URL
	ToURL      string             `bson:"to_url"`      // Destination page URL
	FromDocID  primitive.ObjectID `bson:"from_doc_id"` // Reference to source page document
	ToDocID    primitive.ObjectID `bson:"to_doc_id"`   // Reference to destination page document
	AnchorText string             `bson:"anchor_text"` // Link text for better relevance scoring
	CreatedAt  time.Time          `bson:"created_at"`
}

// PageRank represents the computed score of a URL
// Collection: "pagerank"
type PageRank struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	URL           string             `bson:"url"`
	Score         float64            `bson:"score"`
	InboundCount  int                `bson:"inbound_count"`  // Number of incoming links
	OutboundCount int                `bson:"outbound_count"` // Number of outgoing links
	LastUpdated   time.Time          `bson:"last_updated"`
}
