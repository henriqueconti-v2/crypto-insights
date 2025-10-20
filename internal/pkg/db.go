package pkg

import (
	"database/sql"
	"fmt"

	"crypto-alerts/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn   *sql.DB
	config *config.DatabaseConfig
}

func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao testar conexão com o banco de dados: %w", err)
	}

	return &DB{
		Conn:   conn,
		config: cfg,
	}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
