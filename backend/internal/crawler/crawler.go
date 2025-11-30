// Crawler logic
package crawler

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"searching-engine/internal/storage"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	db       *sql.DB
	visited  map[string]bool
	maxPages int
}

// Initialize a new Crawler
func NewCrawler(db *sql.DB, maxPages int) *Crawler {
	return &Crawler{
		db:       db,
		visited:  make(map[string]bool),
		maxPages: maxPages,
	}
}

// Use queue for BFS crawling
func (c *Crawler) Crawl(seedURL string) error {
	queue := []string{seedURL}
	count := 0

	for len(queue) > 0 && count < c.maxPages {
		currentURL := queue[0]
		queue = queue[1:]

		if c.visited[currentURL] {
			continue
		}

		fmt.Printf("Crawling: %s\n", currentURL)

		page, links, err := c.fetchPage(currentURL)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", currentURL, err)
			continue
		}

		pageID, err := storage.SavePage(c.db, page.URL, page.Title, page.Content)
		if err != nil {
			fmt.Printf("Error saving %s: %v\n", currentURL, err)
			continue
		}

		c.visited[currentURL] = true
		count++

		for _, link := range links {
			absURL := c.resolveURL(currentURL, link)
			if absURL != "" && !c.visited[absURL] {
				queue = append(queue, absURL)
				storage.SaveLink(c.db, pageID, absURL)
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (c *Crawler) fetchPage(pageURL string) (*Page, []string, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	title, content, links := c.parseHTML(string(body))

	page := &Page{
		URL:     pageURL,
		Title:   title,
		Content: content,
	}

	return page, links, nil
}

func (c *Crawler) parseHTML(htmlContent string) (string, string, []string) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", "", nil
	}

	var title string
	var content strings.Builder
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			} else if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			} else if n.Data == "p" || n.Data == "h1" || n.Data == "h2" {
				if n.FirstChild != nil {
					content.WriteString(n.FirstChild.Data + " ")
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(doc)
	return title, content.String(), links
}

func (c *Crawler) resolveURL(base, ref string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	refURL, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	absURL := baseURL.ResolveReference(refURL)

	// Only crawl http/https
	if absURL.Scheme != "http" && absURL.Scheme != "https" {
		return ""
	}

	return absURL.String()
}

type Page struct {
	URL     string
	Title   string
	Content string
}
