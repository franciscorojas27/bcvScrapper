package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type BinanceRate struct {
	SellPrice float64 `json:"sellPrice"`
	BuyPrice  float64 `json:"buyPrice"`
}

type P2PRequest struct {
	Fiat                      string   `json:"fiat"`
	Page                      int      `json:"page"`
	Rows                      int      `json:"rows"`
	TradeType                 string   `json:"tradeType"`
	Asset                     string   `json:"asset"`
	Countries                 []any    `json:"countries"`
	ProMerchantAds            bool     `json:"proMerchantAds"`
	ShieldMerchantAds         bool     `json:"shieldMerchantAds"`
	FilterType                string   `json:"filterType"`
	Periods                   []any    `json:"periods"`
	AdditionalKycVerifyFilter int      `json:"additionalKycVerifyFilter"`
	PublisherType             any      `json:"publisherType"`
	PayTypes                  []string `json:"payTypes"`
	Classifies                []string `json:"classifies"`
	TradedWith                bool     `json:"tradedWith"`
	Followed                  bool     `json:"followed"`
	TransAmount               int      `json:"transAmount"`
}

func GetBinanceRates(tradeType string, amount int, payTypes ...string) (decimal.Decimal, error) {
	url := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	price := decimal.NewFromFloat(0)

	if payTypes == nil {
		payTypes = []string{}
	}

	requestBody := P2PRequest{
		Fiat:                      "VES",
		Page:                      1,
		Rows:                      1,
		TradeType:                 tradeType,
		Asset:                     "USDT",
		Countries:                 []any{},
		ProMerchantAds:            false,
		ShieldMerchantAds:         false,
		PayTypes:                  payTypes,
		FilterType:                "CLASSIC",
		Periods:                   []any{},
		AdditionalKycVerifyFilter: 0,
		PublisherType:             nil,
		TransAmount:               amount,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return price, fmt.Errorf("Error to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return price, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return price, fmt.Errorf("Error to make request to Binance API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return price, fmt.Errorf("Binance returned status: %d", res.StatusCode)
	}

	var response struct {
		Data []struct {
			Adv struct {
				Price string `json:"price"`
			} `json:"adv"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return price, fmt.Errorf("Error to decode Binance API response: %w", err)
	}

	if len(response.Data) == 0 {
		return price, fmt.Errorf("No data received from Binance API")
	}

	price, err = decimal.NewFromString(response.Data[0].Adv.Price)
	if err != nil {
		return price, fmt.Errorf("Error to parse price from Binance API response: %w", err)
	}

	return price, nil
}
