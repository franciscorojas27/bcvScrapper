package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapeBCV() CurrencyRatesData {
	FinalDataRaw := CurrencyRatesData{
		List: []CurrencyRate{},
		Date: "",
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
			FinalDataRaw.Date = t.UTC().Format(time.RFC3339)
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

		priceFloat, err := strconv.ParseFloat(priceNormalized, 64)

		if err != nil {
			fmt.Printf("❌ Error al convertir el precio a float: %v (raw=%q normalized=%q)\n", err, price, priceNormalized)
			return
		}
		if name != "" && price != "" {
			FinalDataRaw.List = append(FinalDataRaw.List, CurrencyRate{
				Symbol:    name,
				Price:     priceFloat,
				ChangePct: 0,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("❌ Error: %v\n", err)
	})

	c.Visit("https://www.bcv.org.ve/")
	return FinalDataRaw
}
