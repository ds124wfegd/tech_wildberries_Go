package service

import (
	"github.com/ds124wfegd/tech_wildberries_Go/internal/database"
)

type Service struct {
}

func NewService(repos *database.Repository) *Service {
	return &Service{}
}
