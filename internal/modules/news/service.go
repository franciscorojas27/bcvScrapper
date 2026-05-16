package news

import (
	newsbank "bcv/internal/platform/providers/news_bank"
	"bcv/internal/platform/scraper"
	"fmt"
	"log/slog"
	"sync"
)

var getNews = scraper.GetNews

var fetchNews = newsbank.FetchNews

var fetchRSSNews = newsbank.FetchRSSNews

var newsEndpoints = []string{
	"https://finanzasdigital.com/wp-json/wp/v2/posts?per_page=5",
	"https://bitacoraeconomica.com/wp-json/wp/v2/posts?per_page=5",
}

func FetchNewsTitles() ([]string, error) {
	type result struct {
		titles []string
		err    error
	}

	results := make(chan result, len(newsEndpoints)+2)
	var wg sync.WaitGroup

	wg.Go(func() {

		titles, err := getNews()
		if err != nil {
			results <- result{err: err}
			return
		}
		results <- result{titles: titles}
	})

	wg.Go(func() {

		titles, err := fetchRSSNews("https://www.criptonoticias.com/feed/")
		if err != nil {
			results <- result{err: err}
			return
		}
		results <- result{titles: titles}
	})

	for _, url := range newsEndpoints {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()

			titles, err := fetchNews(endpoint)
			if err != nil {
				results <- result{err: err}
				return
			}
			results <- result{titles: titles}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allTitles []string
	for item := range results {
		if item.err != nil {
			slog.Error("Error fetching news source", "error", item.err)
			continue
		}
		allTitles = append(allTitles, item.titles...)
	}

	if len(allTitles) == 0 {
		return nil, fmt.Errorf("no news titles fetched from any source")
	}

	return allTitles, nil
}
