package config

import (
	"fmt"
	"os"
	"strconv"
)

type CryptoConfig struct {
	Symbol string `json:"symbol"`

	ThresholdUp1h   float64 `json:"ThresholdUp1hPercent"`
	ThresholdDown1h float64 `json:"ThresholdDown1hPercent"`

	ThresholdUp24h   float64 `json:"ThresholdUp24hPercent"`
	ThresholdDown24h float64 `json:"ThresholdDown24hPercent"`

	ThresholdUp7d   float64 `json:"ThresholdUp7dPercent"`
	ThresholdDown7d float64 `json:"ThresholdDown7dPercent"`

	ThresholdUp30d   float64 `json:"ThresholdUp30dPercent"`
	ThresholdDown30d float64 `json:"ThresholdDown30dPercent"`

	ThresholdUp60d   float64 `json:"ThresholdUp60dPercent"`
	ThresholdDown60d float64 `json:"ThresholdDown60dPercent"`

	ThresholdUp90d   float64 `json:"ThresholdUp90dPercent"`
	ThresholdDown90d float64 `json:"ThresholdDown90dPercent"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Table    string `json:"table"`
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type APIProviderConfig struct {
	Domain string `json:"domain"`
	APIKey string `json:"api_key"`
}

type APIConfig struct {
	CoinGecko     APIProviderConfig `json:"coingecko"`
	CoinMarketCap APIProviderConfig `json:"coinmarketcap"`
}

type Config struct {
	API      APIConfig      `json:"api"`
	Database DatabaseConfig `json:"database"`
	SMTP     SMTPConfig     `json:"smtp"`
}

func LoadConfig() (*Config, error) {
	config := &Config{
		API: APIConfig{
			CoinGecko: APIProviderConfig{
				Domain: os.Getenv("COINGECKO_DOMAIN"),
				APIKey: os.Getenv("COINGECKO_API_KEY"),
			},
			CoinMarketCap: APIProviderConfig{
				Domain: os.Getenv("COINMARKETCAP_DOMAIN"),
				APIKey: os.Getenv("COINMARKETCAP_API_KEY"),
			},
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     parseEnvInt("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			Table:    os.Getenv("DB_TABLE"),
		},
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     parseEnvInt("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func parseEnvInt(key string) int {
	val := os.Getenv(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func validateConfig(config *Config) error {
	if config.API.CoinGecko.Domain == "" {
		return fmt.Errorf("CoinGecko domain is required")
	}
	if config.API.CoinMarketCap.Domain == "" {
		return fmt.Errorf("CoinMarketCap domain is required")
	}
	if config.API.CoinMarketCap.APIKey == "" {
		return fmt.Errorf("CoinMarketCap API key is required (COINMARKETCAP_API_KEY)")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required (DB_HOST)")
	}
	if config.Database.Port == 0 {
		return fmt.Errorf("database port is required (DB_PORT)")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required (DB_USER)")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required (DB_NAME)")
	}
	if config.Database.Table == "" {
		return fmt.Errorf("database table is required (DB_TABLE)")
	}
	if config.SMTP.Host == "" {
		return fmt.Errorf("SMTP host is required (SMTP_HOST)")
	}
	if config.SMTP.Port == 0 {
		return fmt.Errorf("SMTP port is required (SMTP_PORT)")
	}
	if config.SMTP.Username == "" {
		return fmt.Errorf("SMTP username is required (SMTP_USERNAME)")
	}
	if config.SMTP.Password == "" {
		return fmt.Errorf("SMTP password is required (SMTP_PASSWORD)")
	}

	return nil
}
