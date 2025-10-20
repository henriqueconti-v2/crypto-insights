package api

import (
	"crypto-alerts/internal/config"
	"crypto-alerts/internal/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const USDCurrency = "USD"

type CoinMarketCapRepository interface {
	GetCryptoPrices(symbols []string) (map[string]*entity.CryptoCurrency, error)
	GetFearGreedIndex() (*entity.FearGreedIndex, error)
}

type coinMarketCapRepo struct {
	domain string
	apiKey string
	client *http.Client
}

func NewCoinMarketCapRepository(cfg *config.Config) CoinMarketCapRepository {
	return &coinMarketCapRepo{
		domain: cfg.API.CoinMarketCap.Domain,
		apiKey: cfg.API.CoinMarketCap.APIKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *coinMarketCapRepo) GetCryptoPrices(symbols []string) (map[string]*entity.CryptoCurrency, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols provided")
	}

	symbolsStr := strings.Join(symbols, ",")
	params := url.Values{}
	params.Add("symbol", symbolsStr)

	requestURL := fmt.Sprintf("%s/v2/cryptocurrency/quotes/latest?%s", r.domain, params.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to CoinMarketCap API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response entity.CoinMarketCapResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %w", err)
	}

	if response.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("API error: %d - %s", response.Status.ErrorCode, response.Status.ErrorMessage)
	}

	result := make(map[string]*entity.CryptoCurrency)

	for symbol, cryptoDataList := range response.Data {
		if len(cryptoDataList) == 0 {
			continue
		}

		cryptoData := cryptoDataList[0]

		quoteData, ok := cryptoData.Quote[USDCurrency]
		if !ok {
			continue
		}
		quoteData.Name = cryptoData.Name
		result[symbol] = &quoteData
	}

	return result, nil
}

func (r *coinMarketCapRepo) GetFearGreedIndex() (*entity.FearGreedIndex, error) {
	fearGreedURL := fmt.Sprintf("%s/v3/fear-and-greed/latest", r.domain)

	req, err := http.NewRequest("GET", fearGreedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var fearGreedResp entity.FearGreedResponse
	if err := json.Unmarshal(body, &fearGreedResp); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if fearGreedResp.Status.ErrorCode != "0" {
		return nil, fmt.Errorf("API error: %s", fearGreedResp.Status.ErrorMessage)
	}

	return &entity.FearGreedIndex{
		Value:          fearGreedResp.Data.Value,
		Classification: fearGreedResp.Data.ValueClassification,
	}, nil
}
