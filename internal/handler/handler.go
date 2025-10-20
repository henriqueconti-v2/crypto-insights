package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"crypto-alerts/internal/config"
	"crypto-alerts/internal/entity"
	"crypto-alerts/internal/pkg"
	apiRepo "crypto-alerts/internal/repository/api"
	"crypto-alerts/internal/repository/db"
	notifierRepo "crypto-alerts/internal/repository/notifier"
	"crypto-alerts/internal/usecase"
)

type CheckAlertsResponse struct {
	AlertsTriggered int `json:"alerts_triggered"`
}

type API struct {
	config                  *config.Config
	createAlertUseCase      usecase.CreateAlertUseCase
	executeAlertScanUseCase usecase.ExecuteAlertScanUseCase
	db                      *pkg.DB
}

func NewAPI(cfg *config.Config) (*API, error) {
	database, err := pkg.NewDB(&cfg.Database)
	if err != nil {
		return nil, err
	}

	alertRepo := db.NewAlertThresholdRepository(database)
	coinMarketCapRepo := apiRepo.NewCoinMarketCapRepository(cfg)
	coinGeckoRepo := apiRepo.NewCoinGeckoRepository(cfg)
	emailNotifier := notifierRepo.NewEmailNotifier(&cfg.SMTP)

	return &API{
		config:                  cfg,
		createAlertUseCase:      usecase.NewCreateAlertUseCase(alertRepo),
		executeAlertScanUseCase: usecase.NewExecuteAlertScanUseCase(alertRepo, coinMarketCapRepo, coinGeckoRepo, emailNotifier),
		db:                      database,
	}, nil
}

func (api *API) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/crypto_alert_api/execute", api.handleScanExecution)
	mux.HandleFunc("/crypto_alert_api/create", api.handleCreate)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) handleScanExecution(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	alerts, err := api.executeAlertScanUseCase.Execute()
	if err != nil {
		http.Error(w, "Failed to execute alert scan", http.StatusInternalServerError)
		return
	}

	response := CheckAlertsResponse{
		AlertsTriggered: len(alerts),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *API) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var alertThreshold entity.AlertThreshold
	if err := json.NewDecoder(r.Body).Decode(&alertThreshold); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.createAlertUseCase.Execute(&alertThreshold)
	if err != nil {
		log.Printf("Error creating alert: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Configuration saved successfully",
	})
}
