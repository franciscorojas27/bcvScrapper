package domain

import "github.com/shopspring/decimal"

type BinanceRate struct {
	SellPrice decimal.Decimal `json:"sellPrice"`
	BuyPrice  decimal.Decimal `json:"buyPrice"`
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
