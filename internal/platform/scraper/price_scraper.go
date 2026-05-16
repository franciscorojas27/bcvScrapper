package scraper

import (
	"bcv/internal/modules/telegram"
	"bcv/internal/platform/database"
	"bcv/internal/platform/server"
	"bcv/models"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/shopspring/decimal"
)

func scrapeBCV() models.ScrapeReport {
	FinalDataRaw := models.ScrapeReport{
		Rates:   []models.CurrencyRate{},
		BcvDate: time.Time{},
	}

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "www.bcv.org.ve",
		},
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

func ScrapeLatestRates(app *server.App) error {
	data := scrapeBCV()
	if err := database.SaveScrapeReport(app.DB, data); err != nil {
		return fmt.Errorf("Error to save scrape report: %w", err)
	}
	message := telegram.BuildMessage(data)

	if err := telegram.SendMessage(app.Auth, message); err != nil {
		return fmt.Errorf("Error sending message to telegram: %w ", err)
	}
	return nil
}
