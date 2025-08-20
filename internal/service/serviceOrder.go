package service

import (
	"context"
	"errors"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
)

// return the order by id
func (s *Service) GetByUID(ctx context.Context, orderUID string) (*entity.Order, error) {
	if orderUID == "" {
		return nil, errors.New("order_uid is required")
	}

	if order, ok := s.cache.Get(orderUID); ok {
		return order, nil
	}

	order, err := s.repository.GetByUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	if order != nil {
		s.cache.Set(order)
	}
	return order, nil
}

// saves the order, update the cache
func (s *Service) Ingest(ctx context.Context, order *entity.Order) error {

	if order == nil {
		return errors.New("order is nil")
	}
	if order.OrderUID == "" {
		return errors.New("order_uid is required")
	}

	if err := s.repository.Save(ctx, order); err != nil {
		return err
	}
	s.cache.Set(order)
	return nil
}
