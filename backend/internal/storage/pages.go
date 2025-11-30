// Pages table operations
package storage

import (
	"database/sql"
)

func SavePage(db *sql.DB, url, title, content string) (int64, error) {
	result, err := db.Exec(
		"INSERT OR REPLACE INTO pages (url, title, content) VALUES (?, ?, ?)",
		url, title, content,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func SaveLink(db *sql.DB, sourcePageID int64, targetURL string) error {
	var targetPageID int64
	err := db.QueryRow("SELECT id FROM pages WHERE url = ?", targetURL).Scan(&targetPageID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT OR IGNORE INTO links (source_page_id, target_page_id) VALUES (?, ?)",
		sourcePageID, targetPageID,
	)
	return err
}
