package usecase

import (
	"crypto-alerts/internal/entity"
	"crypto-alerts/internal/pkg"
	apiRepo "crypto-alerts/internal/repository/api"
	dbRepo "crypto-alerts/internal/repository/db"
	notifierRepo "crypto-alerts/internal/repository/notifier"
	"log"
)

type ExecuteAlertScanUseCase interface {
	Execute() ([]pkg.AlertMessage, error)
}

type executeAlertScanUseCase struct {
	alertRepo         dbRepo.AlertThresholdRepository
	coinMarketCapRepo apiRepo.CoinMarketCapRepository
	coinGeckoRepo     apiRepo.CoinGeckoRepository
	notifier          notifierRepo.Notifier
}

func NewExecuteAlertScanUseCase(
	alertRepo dbRepo.AlertThresholdRepository,
	coinMarketCapRepo apiRepo.CoinMarketCapRepository,
	coinGeckoRepo apiRepo.CoinGeckoRepository,
	notifier notifierRepo.Notifier,
) ExecuteAlertScanUseCase {
	return &executeAlertScanUseCase{
		alertRepo:         alertRepo,
		coinMarketCapRepo: coinMarketCapRepo,
		coinGeckoRepo:     coinGeckoRepo,
		notifier:          notifier,
	}
}

func (uc *executeAlertScanUseCase) Execute() ([]pkg.AlertMessage, error) {
	thresholds, err := uc.alertRepo.GetAllThresholds()
	if err != nil {
		log.Printf("Error getting thresholds from database: %v", err)
		return nil, err
	}

	if len(thresholds) == 0 {
		log.Println("No thresholds found in database")
		return []pkg.AlertMessage{}, nil
	}

	symbolsMap := make(map[string]bool)
	for _, threshold := range thresholds {
		symbolsMap[threshold.CryptoSymbol] = true
	}

	symbols := make([]string, 0, len(symbolsMap))
	for symbol := range symbolsMap {
		symbols = append(symbols, symbol)
	}

	cryptoData, err := uc.coinMarketCapRepo.GetCryptoPrices(symbols)
	if err != nil {
		log.Printf("Error getting crypto prices: %v", err)
		return nil, err
	}

	fearGreed, err := uc.coinMarketCapRepo.GetFearGreedIndex()
	if err != nil {
		log.Printf("Warning: Failed to get Fear & Greed Index: %v", err)
		fearGreed = nil
	}

	historicalDataMap := make(map[string]*entity.HistoricalPriceData)
	for _, symbol := range symbols {
		historicalData, err := uc.coinGeckoRepo.GetHistoricalPrices(symbol, 90)
		if err != nil {
			log.Printf("Warning: Failed to get historical data for %s: %v", symbol, err)
			continue
		}
		historicalDataMap[symbol] = historicalData
	}

	alerts := uc.processAlerts(thresholds, cryptoData, fearGreed, historicalDataMap)

	log.Printf("Processed %d thresholds and generated %d alerts", len(thresholds), len(alerts))

	return alerts, nil
}

func (uc *executeAlertScanUseCase) processAlerts(
	thresholds []*entity.AlertThreshold,
	cryptoData map[string]*entity.CryptoCurrency,
	fearGreed *entity.FearGreedIndex,
	historicalDataMap map[string]*entity.HistoricalPriceData,
) []pkg.AlertMessage {
	var alerts []pkg.AlertMessage

	for _, threshold := range thresholds {
		data, exists := cryptoData[threshold.CryptoSymbol]
		if !exists {
			continue
		}

		historicalData := historicalDataMap[threshold.CryptoSymbol]

		alertsFound := false

		if threshold.ThresholdUp1hEnabled && threshold.ThresholdUp1hPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "1h", data.PercentChange1h,
				*threshold.ThresholdUp1hPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown1hEnabled && threshold.ThresholdDown1hPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "1h", data.PercentChange1h,
				*threshold.ThresholdDown1hPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.ThresholdUp24hEnabled && threshold.ThresholdUp24hPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "24h", data.PercentChange24h,
				*threshold.ThresholdUp24hPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown24hEnabled && threshold.ThresholdDown24hPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "24h", data.PercentChange24h,
				*threshold.ThresholdDown24hPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.ThresholdUp7dEnabled && threshold.ThresholdUp7dPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "7d", data.PercentChange7d,
				*threshold.ThresholdUp7dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown7dEnabled && threshold.ThresholdDown7dPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "7d", data.PercentChange7d,
				*threshold.ThresholdDown7dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.ThresholdUp30dEnabled && threshold.ThresholdUp30dPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "30d", data.PercentChange30d,
				*threshold.ThresholdUp30dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown30dEnabled && threshold.ThresholdDown30dPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "30d", data.PercentChange30d,
				*threshold.ThresholdDown30dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.ThresholdUp60dEnabled && threshold.ThresholdUp60dPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "60d", data.PercentChange60d,
				*threshold.ThresholdUp60dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown60dEnabled && threshold.ThresholdDown60dPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "60d", data.PercentChange60d,
				*threshold.ThresholdDown60dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.ThresholdUp90dEnabled && threshold.ThresholdUp90dPercent != nil {
			alertsFound = uc.checkUserVarThresholdUp(threshold, data, "90d", data.PercentChange90d,
				*threshold.ThresholdUp90dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.ThresholdDown90dEnabled && threshold.ThresholdDown90dPercent != nil {
			alertsFound = uc.checkUserVarThresholdDown(threshold, data, "90d", data.PercentChange90d,
				*threshold.ThresholdDown90dPercent, fearGreed, historicalData, &alerts) || alertsFound
		}

		if threshold.TargetPriceUpEnabled && threshold.TargetPriceUp != nil {
			alertsFound = uc.checkUserTargetPriceUp(threshold, data, *threshold.TargetPriceUp, fearGreed, historicalData, &alerts) || alertsFound
		}
		if threshold.TargetPriceDownEnabled && threshold.TargetPriceDown != nil {
			alertsFound = uc.checkUserTargetPriceDown(threshold, data, *threshold.TargetPriceDown, fearGreed, historicalData, &alerts) || alertsFound
		}

		if !alertsFound {
			log.Printf("✓ No alerts for user %s - %s - all variations within thresholds",
				threshold.Email, threshold.CryptoSymbol)
		}
	}

	return alerts
}

func (uc *executeAlertScanUseCase) checkUserVarThresholdUp(
	threshold *entity.AlertThreshold,
	data *entity.CryptoCurrency,
	period string,
	variation float64,
	thresholdValue float64,
	fearGreed *entity.FearGreedIndex,
	historicalData *entity.HistoricalPriceData,
	alerts *[]pkg.AlertMessage,
) bool {
	if variation >= thresholdValue {
		alert := pkg.AlertMessage{
			Name:           data.Name,
			Symbol:         threshold.CryptoSymbol,
			Price:          data.Price,
			Volume:         data.Volume24h,
			Period:         period,
			Variation:      variation,
			Threshold:      thresholdValue,
			Direction:      "up",
			HistoricalData: historicalData,
		}

		if fearGreed != nil {
			alert.FearGreedValue = fearGreed.Value
			alert.FearGreedClass = fearGreed.Classification
		}

		*alerts = append(*alerts, alert)

		uc.sendAlertEmailToUser(threshold.Email, alert)

		return true
	}
	return false
}

func (uc *executeAlertScanUseCase) checkUserVarThresholdDown(
	threshold *entity.AlertThreshold,
	data *entity.CryptoCurrency,
	period string,
	variation float64,
	thresholdValue float64,
	fearGreed *entity.FearGreedIndex,
	historicalData *entity.HistoricalPriceData,
	alerts *[]pkg.AlertMessage,
) bool {
	if variation <= thresholdValue {
		alert := pkg.AlertMessage{
			Name:           data.Name,
			Symbol:         threshold.CryptoSymbol,
			Price:          data.Price,
			Volume:         data.Volume24h,
			Period:         period,
			Variation:      variation,
			Threshold:      thresholdValue,
			Direction:      "down",
			HistoricalData: historicalData,
		}

		if fearGreed != nil {
			alert.FearGreedValue = fearGreed.Value
			alert.FearGreedClass = fearGreed.Classification
		}

		*alerts = append(*alerts, alert)

		uc.sendAlertEmailToUser(threshold.Email, alert)

		return true
	}
	return false
}

func (uc *executeAlertScanUseCase) sendAlertEmailToUser(userEmail string, alert pkg.AlertMessage) {
	subject := pkg.FormatEmailSubject(alert)
	body := pkg.FormatEmailBody(alert)

	if err := uc.notifier.SendEmailAlert(userEmail, subject, body); err != nil {
		log.Printf("Failed to send email alert to %s: %v", userEmail, err)
	} else {
		log.Printf("Email alert sent to %s for %s %s %s", userEmail, alert.Symbol, alert.Period, alert.Direction)
	}
}

func (uc *executeAlertScanUseCase) checkUserTargetPriceUp(
	threshold *entity.AlertThreshold,
	data *entity.CryptoCurrency,
	targetPrice float64,
	fearGreed *entity.FearGreedIndex,
	historicalData *entity.HistoricalPriceData,
	alerts *[]pkg.AlertMessage,
) bool {
	if data.Price >= targetPrice {
		alert := pkg.AlertMessage{
			Name:           data.Name,
			Symbol:         threshold.CryptoSymbol,
			Price:          data.Price,
			Volume:         data.Volume24h,
			Period:         "target",
			Variation:      0,
			Threshold:      0,
			Direction:      "up",
			IsTargetPrice:  true,
			TargetPrice:    targetPrice,
			HistoricalData: historicalData,
		}

		if fearGreed != nil {
			alert.FearGreedValue = fearGreed.Value
			alert.FearGreedClass = fearGreed.Classification
		}

		*alerts = append(*alerts, alert)

		uc.sendAlertEmailToUser(threshold.Email, alert)
		return true
	}
	return false
}

func (uc *executeAlertScanUseCase) checkUserTargetPriceDown(
	threshold *entity.AlertThreshold,
	data *entity.CryptoCurrency,
	targetPrice float64,
	fearGreed *entity.FearGreedIndex,
	historicalData *entity.HistoricalPriceData,
	alerts *[]pkg.AlertMessage,
) bool {
	if data.Price <= targetPrice {
		alert := pkg.AlertMessage{
			Name:           data.Name,
			Symbol:         threshold.CryptoSymbol,
			Price:          data.Price,
			Volume:         data.Volume24h,
			Period:         "target",
			Variation:      0,
			Threshold:      0,
			Direction:      "down",
			IsTargetPrice:  true,
			TargetPrice:    targetPrice,
			HistoricalData: historicalData,
		}

		if fearGreed != nil {
			alert.FearGreedValue = fearGreed.Value
			alert.FearGreedClass = fearGreed.Classification
		}

		*alerts = append(*alerts, alert)

		uc.sendAlertEmailToUser(threshold.Email, alert)
		return true
	}
	return false
}
