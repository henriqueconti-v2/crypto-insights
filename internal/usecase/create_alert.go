package usecase

import (
	"crypto-alerts/internal/entity"
	"crypto-alerts/internal/repository/db"
	"fmt"
)

type CreateAlertUseCase interface {
	Execute(alertThreshold *entity.AlertThreshold) error
}

type createAlertUseCase struct {
	alertRepo db.AlertThresholdRepository
}

func NewCreateAlertUseCase(alertRepo db.AlertThresholdRepository) CreateAlertUseCase {
	return &createAlertUseCase{
		alertRepo: alertRepo,
	}
}

func (uc *createAlertUseCase) Execute(alertThreshold *entity.AlertThreshold) error {
	if err := uc.validate(alertThreshold); err != nil {
		return err
	}

	return uc.alertRepo.Create(alertThreshold)
}

func (uc *createAlertUseCase) validate(alertThreshold *entity.AlertThreshold) error {
	if alertThreshold.Email == "" {
		return fmt.Errorf("email is required")
	}

	if alertThreshold.CryptoSymbol == "" {
		return fmt.Errorf("crypto symbol is required")
	}

	if alertThreshold.ThresholdUp1hEnabled && alertThreshold.ThresholdUp1hPercent == nil {
		return fmt.Errorf("threshold up 1h percent is required when enabled")
	}
	if alertThreshold.ThresholdDown1hEnabled && alertThreshold.ThresholdDown1hPercent == nil {
		return fmt.Errorf("threshold down 1h percent is required when enabled")
	}
	if alertThreshold.ThresholdDown1hEnabled && alertThreshold.ThresholdDown1hPercent != nil && *alertThreshold.ThresholdDown1hPercent >= 0 {
		return fmt.Errorf("threshold down 1h percent must be negative")
	}

	if alertThreshold.ThresholdUp24hEnabled && alertThreshold.ThresholdUp24hPercent == nil {
		return fmt.Errorf("threshold up 24h percent is required when enabled")
	}
	if alertThreshold.ThresholdDown24hEnabled && alertThreshold.ThresholdDown24hPercent == nil {
		return fmt.Errorf("threshold down 24h percent is required when enabled")
	}
	if alertThreshold.ThresholdDown24hEnabled && alertThreshold.ThresholdDown24hPercent != nil && *alertThreshold.ThresholdDown24hPercent >= 0 {
		return fmt.Errorf("threshold down 24h percent must be negative")
	}

	if alertThreshold.ThresholdUp7dEnabled && alertThreshold.ThresholdUp7dPercent == nil {
		return fmt.Errorf("threshold up 7d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown7dEnabled && alertThreshold.ThresholdDown7dPercent == nil {
		return fmt.Errorf("threshold down 7d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown7dEnabled && alertThreshold.ThresholdDown7dPercent != nil && *alertThreshold.ThresholdDown7dPercent >= 0 {
		return fmt.Errorf("threshold down 7d percent must be negative")
	}

	if alertThreshold.ThresholdUp30dEnabled && alertThreshold.ThresholdUp30dPercent == nil {
		return fmt.Errorf("threshold up 30d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown30dEnabled && alertThreshold.ThresholdDown30dPercent == nil {
		return fmt.Errorf("threshold down 30d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown30dEnabled && alertThreshold.ThresholdDown30dPercent != nil && *alertThreshold.ThresholdDown30dPercent >= 0 {
		return fmt.Errorf("threshold down 30d percent must be negative")
	}

	if alertThreshold.ThresholdUp60dEnabled && alertThreshold.ThresholdUp60dPercent == nil {
		return fmt.Errorf("threshold up 60d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown60dEnabled && alertThreshold.ThresholdDown60dPercent == nil {
		return fmt.Errorf("threshold down 60d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown60dEnabled && alertThreshold.ThresholdDown60dPercent != nil && *alertThreshold.ThresholdDown60dPercent >= 0 {
		return fmt.Errorf("threshold down 60d percent must be negative")
	}

	if alertThreshold.ThresholdUp90dEnabled && alertThreshold.ThresholdUp90dPercent == nil {
		return fmt.Errorf("threshold up 90d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown90dEnabled && alertThreshold.ThresholdDown90dPercent == nil {
		return fmt.Errorf("threshold down 90d percent is required when enabled")
	}
	if alertThreshold.ThresholdDown90dEnabled && alertThreshold.ThresholdDown90dPercent != nil && *alertThreshold.ThresholdDown90dPercent >= 0 {
		return fmt.Errorf("threshold down 90d percent must be negative")
	}

	if alertThreshold.TargetPriceUpEnabled && alertThreshold.TargetPriceUp == nil {
		return fmt.Errorf("target price up is required when enabled")
	}
	if alertThreshold.TargetPriceUpEnabled && alertThreshold.TargetPriceUp != nil && *alertThreshold.TargetPriceUp <= 0 {
		return fmt.Errorf("target price up must be positive")
	}

	if alertThreshold.TargetPriceDownEnabled && alertThreshold.TargetPriceDown == nil {
		return fmt.Errorf("target price down is required when enabled")
	}
	if alertThreshold.TargetPriceDownEnabled && alertThreshold.TargetPriceDown != nil && *alertThreshold.TargetPriceDown <= 0 {
		return fmt.Errorf("target price down must be positive")
	}

	return nil
}
