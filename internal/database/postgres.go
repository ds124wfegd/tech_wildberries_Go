package database

import (
	"fmt"

	"github.com/ds124wfegd/tech_wildberries_Go/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// create DB entity
func NewPostgresDB(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(cfg.PgDriver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode, cfg.Password))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
