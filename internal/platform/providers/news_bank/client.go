package newsbank

import (
	"bcv/internal/platform/scraper"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"log/slog"
)

type WPPost struct {
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
}

func fetchNews(url string) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api responded with status: %d", resp.StatusCode)
	}

	var posts []WPPost
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, err
	}

	var titles []string
	for _, post := range posts {
		titles = append(titles, post.Title.Rendered)
	}

	return titles, nil
}

func FetchNewsTitles() ([]string, error) {
	data, err := scraper.GetNews()
	if err != nil {
		slog.Error("Error fetching news", "error", err)
	}
	endpoints := []string{
		"https://finanzasdigital.com/wp-json/wp/v2/posts?per_page=5",
		"https://bitacoraeconomica.com/wp-json/wp/v2/posts?per_page=5",
	}

	var allTitles []string

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range endpoints {
		wg.Add(1)
		go func(enpoint string) {
			defer wg.Done()

			titles, err := fetchNews(url)
			if err != nil {
				slog.Error("Error fetching news from endpoint", "url", url, "error", err)
				return
			}
			mu.Lock()
			allTitles = append(allTitles, titles...)
			mu.Unlock()
		}(url)
	}

	if len(allTitles) == 0 {
		return nil, fmt.Errorf("no news titles fetched from any source")
	}

	data = append(data, allTitles...)

	wg.Wait()

	return data, nil
}
