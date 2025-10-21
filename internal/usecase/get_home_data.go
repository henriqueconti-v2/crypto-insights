package usecase

import (
	"crypto-alerts/internal/entity"
	"crypto-alerts/internal/repository/api"
	"fmt"
	"log"
)

type GetHomeDataUseCase interface {
	Execute() (*entity.HomeData, error)
}

type getHomeDataUseCase struct {
	coinMarketCapRepo api.CoinMarketCapRepository
}

func NewGetHomeDataUseCase(coinMarketCapRepo api.CoinMarketCapRepository) GetHomeDataUseCase {
	return &getHomeDataUseCase{
		coinMarketCapRepo: coinMarketCapRepo,
	}
}

func (uc *getHomeDataUseCase) Execute() (*entity.HomeData, error) {
	fearGreedIndex, err := uc.coinMarketCapRepo.GetFearGreedIndex()
	if err != nil {
		log.Printf("Error fetching Fear & Greed Index: %v", err)
		return nil, fmt.Errorf("failed to fetch Fear & Greed Index: %w", err)
	}

	label := uc.getFearGreedLabel(fearGreedIndex.Value)

	return &entity.HomeData{
		FearGreedIndex: entity.FearGreedIndexData{
			Value:          fearGreedIndex.Value,
			Classification: fearGreedIndex.Classification,
			Label:          label,
		},
	}, nil
}

func (uc *getHomeDataUseCase) getFearGreedLabel(value int) string {
	if value <= 25 {
		return "Medo Extremo"
	}
	if value <= 45 {
		return "Medo"
	}
	if value <= 55 {
		return "Neutro"
	}
	if value <= 75 {
		return "Ganância"
	}
	return "Ganância Extrema"
}
