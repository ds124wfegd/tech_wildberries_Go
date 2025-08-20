package database

import (
	"context"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Order-repository interface
type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	GetByUID(ctx context.Context, orderUID string) (*entity.Order, error)
	GetRecentUIDs(ctx context.Context, limit int) ([]string, error)
}

type Repository struct {
	OrderRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		OrderRepository: NewOrderPostgres(db),
	}
}
