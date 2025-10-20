package api

import (
	"crypto-alerts/internal/config"
	"crypto-alerts/internal/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var coinGeckoIDMap = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"SOL":   "solana",
	"BNB":   "binancecoin",
	"XRP":   "ripple",
	"ADA":   "cardano",
	"DOGE":  "dogecoin",
	"MATIC": "matic-network",
	"DOT":   "polkadot",
	"AVAX":  "avalanche-2",
}

type CoinGeckoRepository interface {
	GetHistoricalPrices(symbol string, days int) (*entity.HistoricalPriceData, error)
}

type coinGeckoRepo struct {
	domain string
	apiKey string
	client *http.Client
}

func NewCoinGeckoRepository(cfg *config.Config) CoinGeckoRepository {
	return &coinGeckoRepo{
		domain: cfg.API.CoinGecko.Domain,
		apiKey: cfg.API.CoinGecko.APIKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *coinGeckoRepo) GetHistoricalPrices(symbol string, days int) (*entity.HistoricalPriceData, error) {
	coinID, exists := coinGeckoIDMap[strings.ToUpper(symbol)]
	if !exists {
		return nil, fmt.Errorf("símbolo %s não suportado pela CoinGecko", symbol)
	}

	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=%d&interval=daily", r.domain, coinID, days)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if r.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", r.apiKey)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados históricos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na API CoinGecko (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResponse entity.CoinGeckoMarketChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	historicalData := &entity.HistoricalPriceData{
		Symbol:    symbol,
		Prices:    make([]entity.PriceHistoryPoint, 0, len(apiResponse.Prices)),
		Volumes:   make([]entity.VolumePoint, 0, len(apiResponse.TotalVolumes)),
		DaysCount: len(apiResponse.Prices),
	}

	var priceSum float64
	minPrice := float64(0)
	maxPrice := float64(0)

	for i, priceData := range apiResponse.Prices {
		if len(priceData) < 2 {
			continue
		}

		timestamp := int64(priceData[0])
		price := priceData[1]

		historicalData.Prices = append(historicalData.Prices, entity.PriceHistoryPoint{
			Timestamp: timestamp,
			Price:     price,
		})

		priceSum += price

		if i == 0 || price < minPrice {
			minPrice = price
		}
		if i == 0 || price > maxPrice {
			maxPrice = price
		}
	}

	if len(historicalData.Prices) > 0 {
		historicalData.MinPrice = minPrice
		historicalData.MaxPrice = maxPrice
		historicalData.AvgPrice = priceSum / float64(len(historicalData.Prices))
	}

	var volumeSum float64
	minVolume := float64(0)
	maxVolume := float64(0)

	for i, volumeData := range apiResponse.TotalVolumes {
		if len(volumeData) < 2 {
			continue
		}

		timestamp := int64(volumeData[0])
		volume := volumeData[1]

		historicalData.Volumes = append(historicalData.Volumes, entity.VolumePoint{
			Timestamp: timestamp,
			Volume:    volume,
		})

		volumeSum += volume

		if i == 0 || volume < minVolume {
			minVolume = volume
		}
		if i == 0 || volume > maxVolume {
			maxVolume = volume
		}
	}

	if len(historicalData.Volumes) > 0 {
		historicalData.MinVolume = minVolume
		historicalData.MaxVolume = maxVolume
		historicalData.AvgVolume = volumeSum / float64(len(historicalData.Volumes))
	}

	return historicalData, nil
}
