package service

import (
	"context"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/repository"
)

// ProductService contains the business logic for product operations.
// It depends only on the ProductRepository interface, keeping it storage-agnostic.
type ProductService struct {
	products repository.ProductRepository
}

// NewProductService constructs a ProductService with its required dependencies.
func NewProductService(products repository.ProductRepository) *ProductService {
	return &ProductService{products: products}
}

// GetAllProducts returns the full product catalogue.
func (s *ProductService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.products.FindAll(ctx)
}

// GetProduct returns a single product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.products.FindByID(ctx, id)
}
