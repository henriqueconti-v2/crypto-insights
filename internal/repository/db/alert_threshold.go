package db

import (
	"crypto-alerts/internal/entity"
	"crypto-alerts/internal/pkg"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type AlertThresholdRepository interface {
	Create(threshold *entity.AlertThreshold) error
	GetAllThresholds() ([]*entity.AlertThreshold, error)
}

type AlertThresholdPostgres struct {
	db *pkg.DB
}

func NewAlertThresholdRepository(db *pkg.DB) AlertThresholdRepository {
	return &AlertThresholdPostgres{db: db}
}

func (r *AlertThresholdPostgres) Create(threshold *entity.AlertThreshold) error {
	rand.Seed(time.Now().UnixNano())

	id := rand.Intn(1000000000)

	sql := `
		INSERT INTO user_crypto_thresholds (
			id, 
			email, 
			crypto_symbol,
			
			threshold_up_1h_percent, 
			threshold_up_1h_enabled, 
			threshold_down_1h_percent, 
			threshold_down_1h_enabled,
			
			threshold_up_24h_percent, 
			threshold_up_24h_enabled, 
			threshold_down_24h_percent, 
			threshold_down_24h_enabled,
			
			threshold_up_7d_percent, 
			threshold_up_7d_enabled, 
			threshold_down_7d_percent, 
			threshold_down_7d_enabled,
			
			threshold_up_30d_percent, 
			threshold_up_30d_enabled, 
			threshold_down_30d_percent, 
			threshold_down_30d_enabled,
			
			threshold_up_60d_percent, 
			threshold_up_60d_enabled, 
			threshold_down_60d_percent, 
			threshold_down_60d_enabled,
			
			threshold_up_90d_percent, 
			threshold_up_90d_enabled, 
			threshold_down_90d_percent, 
			threshold_down_90d_enabled,
			
			target_price_up,
			target_price_up_enabled,
			target_price_down,
			target_price_down_enabled,
			
			created_at
		) VALUES (
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8, $9, $10, $11, 
			$12, $13, $14, $15, 
			$16, $17, $18, $19, 
			$20, $21, $22, $23, 
			$24, $25, $26, $27, 
			$28, $29, $30, $31,
			$32
		)
	`

	var up1hPercent, down1hPercent interface{}
	var up24hPercent, down24hPercent interface{}
	var up7dPercent, down7dPercent interface{}
	var up30dPercent, down30dPercent interface{}
	var up60dPercent, down60dPercent interface{}
	var up90dPercent, down90dPercent interface{}
	var targetPriceUp, targetPriceDown interface{}

	if threshold.ThresholdUp1hPercent != nil {
		up1hPercent = *threshold.ThresholdUp1hPercent
	}
	if threshold.ThresholdDown1hPercent != nil {
		down1hPercent = *threshold.ThresholdDown1hPercent
	}
	if threshold.ThresholdUp24hPercent != nil {
		up24hPercent = *threshold.ThresholdUp24hPercent
	}
	if threshold.ThresholdDown24hPercent != nil {
		down24hPercent = *threshold.ThresholdDown24hPercent
	}
	if threshold.ThresholdUp7dPercent != nil {
		up7dPercent = *threshold.ThresholdUp7dPercent
	}
	if threshold.ThresholdDown7dPercent != nil {
		down7dPercent = *threshold.ThresholdDown7dPercent
	}
	if threshold.ThresholdUp30dPercent != nil {
		up30dPercent = *threshold.ThresholdUp30dPercent
	}
	if threshold.ThresholdDown30dPercent != nil {
		down30dPercent = *threshold.ThresholdDown30dPercent
	}
	if threshold.ThresholdUp60dPercent != nil {
		up60dPercent = *threshold.ThresholdUp60dPercent
	}
	if threshold.ThresholdDown60dPercent != nil {
		down60dPercent = *threshold.ThresholdDown60dPercent
	}
	if threshold.ThresholdUp90dPercent != nil {
		up90dPercent = *threshold.ThresholdUp90dPercent
	}
	if threshold.ThresholdDown90dPercent != nil {
		down90dPercent = *threshold.ThresholdDown90dPercent
	}
	if threshold.TargetPriceUp != nil {
		targetPriceUp = *threshold.TargetPriceUp
	}
	if threshold.TargetPriceDown != nil {
		targetPriceDown = *threshold.TargetPriceDown
	}

	_, err := r.db.Conn.Exec(
		sql,
		id,
		threshold.Email,
		threshold.CryptoSymbol,

		up1hPercent,
		threshold.ThresholdUp1hEnabled,
		down1hPercent,
		threshold.ThresholdDown1hEnabled,

		up24hPercent,
		threshold.ThresholdUp24hEnabled,
		down24hPercent,
		threshold.ThresholdDown24hEnabled,

		up7dPercent, // Threshold de alta 7d
		threshold.ThresholdUp7dEnabled,
		down7dPercent, // Threshold de baixa 7d
		threshold.ThresholdDown7dEnabled,

		up30dPercent, // Threshold de alta 30d
		threshold.ThresholdUp30dEnabled,
		down30dPercent,
		threshold.ThresholdDown30dEnabled,

		up60dPercent,
		threshold.ThresholdUp60dEnabled,
		down60dPercent,
		threshold.ThresholdDown60dEnabled,

		up90dPercent,
		threshold.ThresholdUp90dEnabled,
		down90dPercent,
		threshold.ThresholdDown90dEnabled,

		targetPriceUp,
		threshold.TargetPriceUpEnabled,
		targetPriceDown,
		threshold.TargetPriceDownEnabled,

		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar threshold no banco de dados: %w", err)
	}

	log.Printf("Threshold salvo com sucesso no banco de dados com ID: %d", id)
	return nil
}

func (r *AlertThresholdPostgres) GetAllThresholds() ([]*entity.AlertThreshold, error) {
	sql := `
		SELECT 
			email, 
			crypto_symbol,
			
			threshold_up_1h_percent, 
			threshold_up_1h_enabled, 
			threshold_down_1h_percent, 
			threshold_down_1h_enabled,
			
			threshold_up_24h_percent, 
			threshold_up_24h_enabled, 
			threshold_down_24h_percent, 
			threshold_down_24h_enabled,
			
			threshold_up_7d_percent, 
			threshold_up_7d_enabled, 
			threshold_down_7d_percent, 
			threshold_down_7d_enabled,
			
			threshold_up_30d_percent, 
			threshold_up_30d_enabled, 
			threshold_down_30d_percent, 
			threshold_down_30d_enabled,
			
			threshold_up_60d_percent, 
			threshold_up_60d_enabled, 
			threshold_down_60d_percent, 
			threshold_down_60d_enabled,
			
			threshold_up_90d_percent, 
			threshold_up_90d_enabled, 
			threshold_down_90d_percent, 
			threshold_down_90d_enabled,
			
			target_price_up,
			target_price_up_enabled,
			target_price_down,
			target_price_down_enabled
		FROM user_crypto_thresholds
		ORDER BY crypto_symbol, email
	`

	rows, err := r.db.Conn.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar thresholds do banco de dados: %w", err)
	}
	defer rows.Close()

	var thresholds []*entity.AlertThreshold

	for rows.Next() {
		threshold := &entity.AlertThreshold{}
		
		var up1hPercent, down1hPercent *float64
		var up24hPercent, down24hPercent *float64
		var up7dPercent, down7dPercent *float64
		var up30dPercent, down30dPercent *float64
		var up60dPercent, down60dPercent *float64
		var up90dPercent, down90dPercent *float64
		var targetPriceUp, targetPriceDown *float64

		err := rows.Scan(
			&threshold.Email,
			&threshold.CryptoSymbol,

			&up1hPercent,
			&threshold.ThresholdUp1hEnabled,
			&down1hPercent,
			&threshold.ThresholdDown1hEnabled,

			&up24hPercent,
			&threshold.ThresholdUp24hEnabled,
			&down24hPercent,
			&threshold.ThresholdDown24hEnabled,

			&up7dPercent,
			&threshold.ThresholdUp7dEnabled,
			&down7dPercent,
			&threshold.ThresholdDown7dEnabled,

			&up30dPercent,
			&threshold.ThresholdUp30dEnabled,
			&down30dPercent,
			&threshold.ThresholdDown30dEnabled,

			&up60dPercent,
			&threshold.ThresholdUp60dEnabled,
			&down60dPercent,
			&threshold.ThresholdDown60dEnabled,

			&up90dPercent,
			&threshold.ThresholdUp90dEnabled,
			&down90dPercent,
			&threshold.ThresholdDown90dEnabled,

			&targetPriceUp,
			&threshold.TargetPriceUpEnabled,
			&targetPriceDown,
			&threshold.TargetPriceDownEnabled,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan dos thresholds: %w", err)
		}

		threshold.ThresholdUp1hPercent = up1hPercent
		threshold.ThresholdDown1hPercent = down1hPercent
		threshold.ThresholdUp24hPercent = up24hPercent
		threshold.ThresholdDown24hPercent = down24hPercent
		threshold.ThresholdUp7dPercent = up7dPercent
		threshold.ThresholdDown7dPercent = down7dPercent
		threshold.ThresholdUp30dPercent = up30dPercent
		threshold.ThresholdDown30dPercent = down30dPercent
		threshold.ThresholdUp60dPercent = up60dPercent
		threshold.ThresholdDown60dPercent = down60dPercent
		threshold.ThresholdUp90dPercent = up90dPercent
		threshold.ThresholdDown90dPercent = down90dPercent
		threshold.TargetPriceUp = targetPriceUp
		threshold.TargetPriceDown = targetPriceDown

		thresholds = append(thresholds, threshold)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	log.Printf("Carregados %d thresholds do banco de dados", len(thresholds))
	return thresholds, nil
}
