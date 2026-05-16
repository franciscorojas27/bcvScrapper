package pricebank

import (
	newsbank "bcv/internal/platform/providers/news_bank"
	"bcv/internal/platform/scraper"
	"fmt"
	"log/slog"
	"sync"
)

var getNews = scraper.GetNews

var fetchNews = newsbank.FetchNews

var newsEndpoints = []string{
	"https://finanzasdigital.com/wp-json/wp/v2/posts?per_page=5",
	"https://bitacoraeconomica.com/wp-json/wp/v2/posts?per_page=5",
}

func FetchNewsTitles() ([]string, error) {
	data, err := getNews()
	if err != nil {
		slog.Error("Error fetching news", "error", err)
	}

	var allTitles []string

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range newsEndpoints {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()

			titles, err := fetchNews(endpoint)
			if err != nil {
				slog.Error("Error fetching news from endpoint", "url", endpoint, "error", err)
				return
			}
			mu.Lock()
			allTitles = append(allTitles, titles...)
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	if len(allTitles) == 0 {
		return nil, fmt.Errorf("no news titles fetched from any source")
	}

	data = append(data, allTitles...)

	return data, nil
}
