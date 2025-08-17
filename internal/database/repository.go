package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
