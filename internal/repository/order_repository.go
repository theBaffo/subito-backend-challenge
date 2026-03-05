package repository

import (
	"context"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
)

// OrderRepository defines the contract for order persistence.
type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	FindByID(ctx context.Context, id string) (*domain.Order, error)
}
