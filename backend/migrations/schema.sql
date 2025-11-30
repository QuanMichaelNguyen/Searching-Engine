-- Database schema
-- Pages table: stores crawled web pages
CREATE TABLE IF NOT EXISTS pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    content TEXT,
    crawled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pagerank REAL DEFAULT 0.0
);

-- Links table: stores link graph for PageRank
CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_page_id INTEGER NOT NULL,
    target_page_id INTEGER NOT NULL,
    FOREIGN KEY (source_page_id) REFERENCES pages(id),
    FOREIGN KEY (target_page_id) REFERENCES pages(id),
    UNIQUE(source_page_id, target_page_id)
);

-- Inverted index: maps terms to pages
CREATE TABLE IF NOT EXISTS inverted_index (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    term TEXT NOT NULL,
    page_id INTEGER NOT NULL,
    frequency INTEGER DEFAULT 1,
    FOREIGN KEY (page_id) REFERENCES pages(id),
    UNIQUE(term, page_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pages_url ON pages(url);
CREATE INDEX IF NOT EXISTS idx_pages_pagerank ON pages(pagerank DESC);
CREATE INDEX IF NOT EXISTS idx_links_source ON links(source_page_id);
CREATE INDEX IF NOT EXISTS idx_links_target ON links(target_page_id);
CREATE INDEX IF NOT EXISTS idx_inverted_term ON inverted_index(term);
CREATE INDEX IF NOT EXISTS idx_inverted_page ON inverted_index(page_id);