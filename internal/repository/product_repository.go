package repository

import (
	"context"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
)

// ProductRepository defines the contract for product persistence.
type ProductRepository interface {
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
}
