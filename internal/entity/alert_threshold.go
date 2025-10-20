package entity

type AlertThreshold struct {
	Email        string `json:"email"`
	CryptoSymbol string `json:"crypto_symbol"`

	// 1h thresholds
	ThresholdUp1hPercent   *float64 `json:"threshold_up_1h_percent"`
	ThresholdUp1hEnabled   bool     `json:"threshold_up_1h_enabled"`
	ThresholdDown1hPercent *float64 `json:"threshold_down_1h_percent"`
	ThresholdDown1hEnabled bool     `json:"threshold_down_1h_enabled"`

	// 24h thresholds
	ThresholdUp24hPercent   *float64 `json:"threshold_up_24h_percent"`
	ThresholdUp24hEnabled   bool     `json:"threshold_up_24h_enabled"`
	ThresholdDown24hPercent *float64 `json:"threshold_down_24h_percent"`
	ThresholdDown24hEnabled bool     `json:"threshold_down_24h_enabled"`

	// 7d thresholds
	ThresholdUp7dPercent   *float64 `json:"threshold_up_7d_percent"`
	ThresholdUp7dEnabled   bool     `json:"threshold_up_7d_enabled"`
	ThresholdDown7dPercent *float64 `json:"threshold_down_7d_percent"`
	ThresholdDown7dEnabled bool     `json:"threshold_down_7d_enabled"`

	// 30d thresholds
	ThresholdUp30dPercent   *float64 `json:"threshold_up_30d_percent"`
	ThresholdUp30dEnabled   bool     `json:"threshold_up_30d_enabled"`
	ThresholdDown30dPercent *float64 `json:"threshold_down_30d_percent"`
	ThresholdDown30dEnabled bool     `json:"threshold_down_30d_enabled"`

	// 60d thresholds
	ThresholdUp60dPercent   *float64 `json:"threshold_up_60d_percent"`
	ThresholdUp60dEnabled   bool     `json:"threshold_up_60d_enabled"`
	ThresholdDown60dPercent *float64 `json:"threshold_down_60d_percent"`
	ThresholdDown60dEnabled bool     `json:"threshold_down_60d_enabled"`

	// 90d thresholds
	ThresholdUp90dPercent   *float64 `json:"threshold_up_90d_percent"`
	ThresholdUp90dEnabled   bool     `json:"threshold_up_90d_enabled"`
	ThresholdDown90dPercent *float64 `json:"threshold_down_90d_percent"`
	ThresholdDown90dEnabled bool     `json:"threshold_down_90d_enabled"`

	// Target price thresholds
	TargetPriceUp      *float64 `json:"target_price_up"`
	TargetPriceUpEnabled   bool     `json:"target_price_up_enabled"`
	TargetPriceDown    *float64 `json:"target_price_down"`
	TargetPriceDownEnabled bool     `json:"target_price_down_enabled"`
}
