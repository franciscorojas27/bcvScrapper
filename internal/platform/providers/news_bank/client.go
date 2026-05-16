package newsbank

import (
	"encoding/json"
	"encoding/xml"
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

type RSSFeed struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchRSSNews(url string) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rss responded with status: %d", resp.StatusCode)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	var titles []string
	for i, item := range feed.Channel.Items {
		if i >= 5 {
			break
		}
		titles = append(titles, item.Title)
	}

	return titles, nil
}
