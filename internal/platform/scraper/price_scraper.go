package scraper

import (
	"bcv/internal/domain"
	"bcv/internal/modules/telegram"
	"bcv/internal/platform/database"
	"bcv/models"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func scrapeBCV() models.ScrapeReport {
	FinalDataRaw := models.ScrapeReport{
		Rates:   []models.CurrencyRate{},
		BcvDate: time.Time{},
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	c.SetRequestTimeout(15 * time.Second)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "www.bcv.org.ve",
			MinVersion:         tls.VersionTLS12,
		},
	}

	proxyURL, err := url.Parse("http://194.180.188.100:8080")
	if err == nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	c.WithTransport(transport)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	c.OnHTML(".pull-right.dinpro.center .date-display-single", func(e *colly.HTMLElement) {
		dateRaw := e.Attr("content")
		if dateRaw == "" {
			fmt.Println("❌ No se encontró la fecha en el HTML")
			return
		}
		t, err := time.Parse(time.RFC3339, dateRaw)
		if err != nil {
			fmt.Printf("❌ Error al parsear la fecha: %v\n", err)
		} else {
			FinalDataRaw.BcvDate = t.UTC()
		}
	})
	c.OnHTML("#euro, #yuan, #lira, #rublo, #dolar", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("span"))
		price := strings.TrimSpace(e.ChildText(".centrado strong"))

		priceNormalized := price
		if strings.Contains(price, ",") && strings.Contains(price, ".") {
			priceNormalized = strings.ReplaceAll(priceNormalized, ".", "")
		}
		priceNormalized = strings.ReplaceAll(priceNormalized, ",", ".")
		priceNormalized = strings.ReplaceAll(priceNormalized, " ", "")
		priceNormalized = strings.ReplaceAll(priceNormalized, "\u00A0", "")

		priceFloat, err := decimal.NewFromString(priceNormalized)

		if err != nil {
			fmt.Printf("❌ Error al convertir el precio a float: %v (raw=%q normalized=%q)\n", err, price, priceNormalized)
			return
		}
		if name != "" && price != "" {
			FinalDataRaw.Rates = append(FinalDataRaw.Rates, models.CurrencyRate{
				Symbol:    name,
				Price:     priceFloat,
				ChangePct: decimal.Zero,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("❌ Error: %v\n", err)
	})

	c.Visit("https://www.bcv.org.ve/")
	return FinalDataRaw
}

func ScrapeLatestRates(db *gorm.DB, auth domain.AuthTelegram) error {
	data := scrapeBCV()
	saved, err := database.SaveScrapeReport(db, data)
	if err != nil {
		return fmt.Errorf("Error to save scrape report: %w", err)
	}
	if !saved {
		return nil
	}

	message := telegram.BuildMessage(data)
	if err := telegram.SendMessage(auth, message); err != nil {
		return fmt.Errorf("Error sending message to telegram: %w ", err)
	}

	return nil
}
