package newsbank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WPPost struct {
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
}

func FetchNews(url string) ([]string, error) {
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

