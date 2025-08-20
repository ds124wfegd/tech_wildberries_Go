package service

import (
	"context"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/database"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
)

// Order-Cache interface
type OrderCache interface {
	Get(orderUID string) (*entity.Order, bool)
	Set(order *entity.Order)
	Load(orders []*entity.Order)
}

// Order-Service interface
type OrderService interface {
	GetByUID(ctx context.Context, orderUID string) (*entity.Order, error)
	Ingest(ctx context.Context, order *entity.Order) error
}

type Service struct {
	repository database.OrderRepository
	cache      OrderCache
}

func NewService(repository database.OrderRepository, cache OrderCache) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
	}
}
