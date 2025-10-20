package entity

import (
	"time"
)

type CoinMarketCapResponse struct {
	Status Status                        `json:"status"`
	Data   map[string][]CryptoDataDetail `json:"data"`
}

type Status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type CryptoDataDetail struct {
	Name  string                    `json:"name"`
	Quote map[string]CryptoCurrency `json:"quote"`
}

type CryptoCurrency struct {
	Price              float64   `json:"price"`
	Volume24h          float64   `json:"volume_24h"`
	VolumeChange24h    float64   `json:"volume_change_24h"`
	PercentChange1h    float64   `json:"percent_change_1h"`
	PercentChange24h   float64   `json:"percent_change_24h"`
	PercentChange7d    float64   `json:"percent_change_7d"`
	PercentChange30d   float64   `json:"percent_change_30d"`
	PercentChange60d   float64   `json:"percent_change_60d"`
	PercentChange90d   float64   `json:"percent_change_90d"`
	MarketCap          float64   `json:"market_cap"`
	MarketCapDominance float64   `json:"market_cap_dominance"`
	Name               string    `json:"name"`
	LastUpdated        time.Time `json:"last_updated"`
}
