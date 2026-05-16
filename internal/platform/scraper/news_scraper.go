package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func GetNews() ([]string, error) {
	var noticias []string

	c := colly.NewCollector(
		colly.AllowedDomains("www.bancaynegocios.com", "bancaynegocios.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(10 * time.Second)

	c.OnHTML("section.recomendaciones_del_editor h2.post-title a, .articulos_mas_leidos h2.post-title a", func(e *colly.HTMLElement) {
		if len(noticias) < 6 {
			titulo := strings.TrimSpace(e.Text)
			if titulo != "" {
				noticias = append(noticias, titulo)
			}
		}
	})

	var errScraping error
	c.OnError(func(r *colly.Response, err error) {
		errScraping = fmt.Errorf("error en scraper (status %d): %w", r.StatusCode, err)
	})

	err := c.Visit("https://www.bancaynegocios.com/")
	if err != nil {
		return nil, err
	}
	if errScraping != nil {
		return nil, errScraping
	}

	return noticias, nil
}
