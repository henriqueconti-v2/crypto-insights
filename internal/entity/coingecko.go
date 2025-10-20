package entity

type CoinGeckoMarketChartResponse struct {
	Prices       [][]float64 `json:"prices"`
	MarketCaps   [][]float64 `json:"market_caps"`
	TotalVolumes [][]float64 `json:"total_volumes"`
}

type PriceHistoryPoint struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

type VolumePoint struct {
	Timestamp int64   `json:"timestamp"`
	Volume    float64 `json:"volume"`
}

type HistoricalPriceData struct {
	Symbol    string              `json:"symbol"`
	Prices    []PriceHistoryPoint `json:"prices"`
	Volumes   []VolumePoint       `json:"volumes"`
	MinPrice  float64             `json:"min_price"`
	MaxPrice  float64             `json:"max_price"`
	AvgPrice  float64             `json:"avg_price"`
	MinVolume float64             `json:"min_volume"`
	MaxVolume float64             `json:"max_volume"`
	AvgVolume float64             `json:"avg_volume"`
	DaysCount int                 `json:"days_count"`
}
