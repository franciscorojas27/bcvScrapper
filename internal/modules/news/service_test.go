package news

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchNewsTitlesAggregatesResults(t *testing.T) {
	originalGetNews := getNews
	originalFetchNews := fetchNews
	originalFetchRSSNews := fetchRSSNews
	originalEndpoints := newsEndpoints

	defer func() {
		getNews = originalGetNews
		fetchNews = originalFetchNews
		fetchRSSNews = originalFetchRSSNews
		newsEndpoints = originalEndpoints
	}()

	getNews = func() ([]string, error) {
		return []string{"base"}, nil
	}
	fetchRSSNews = func(url string) ([]string, error) {
		return []string{"rss"}, nil
	}
	newsEndpoints = []string{"one", "two"}
	fetchNews = func(endpoint string) ([]string, error) {
		switch endpoint {
		case "one":
			return []string{"a", "b"}, nil
		case "two":
			return []string{"c"}, nil
		default:
			return nil, fmt.Errorf("unexpected endpoint %s", endpoint)
		}
	}

	got, err := FetchNewsTitles()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"base", "rss", "a", "b", "c"}, got)
}

func TestFetchNewsTitlesReturnsErrorWhenNoTitlesArrive(t *testing.T) {
	originalGetNews := getNews
	originalFetchNews := fetchNews
	originalFetchRSSNews := fetchRSSNews
	originalEndpoints := newsEndpoints

	defer func() {
		getNews = originalGetNews
		fetchNews = originalFetchNews
		fetchRSSNews = originalFetchRSSNews
		newsEndpoints = originalEndpoints
	}()

	getNews = func() ([]string, error) {
		return nil, nil
	}
	fetchRSSNews = func(url string) ([]string, error) {
		return nil, fmt.Errorf("failed rss")
	}
	newsEndpoints = []string{"one", "two"}
	fetchNews = func(endpoint string) ([]string, error) {
		return nil, fmt.Errorf("failed for %s", endpoint)
	}

	got, err := FetchNewsTitles()
	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no news titles fetched")
}
