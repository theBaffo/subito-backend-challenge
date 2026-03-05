package memory

import (
	"context"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/theBaffo/subito-backend-challenge/internal/domain"
)

// ProductStore is a thread-safe, in-memory implementation of ProductRepository.
// It is pre-seeded with sample products for demonstration purposes.
type ProductStore struct {
	mu       sync.RWMutex
	products map[string]domain.Product
}

// NewProductStore creates a ProductStore pre-populated with sample products.
func NewProductStore() *ProductStore {
	s := &ProductStore{
		products: make(map[string]domain.Product),
	}
	s.seed()
	return s
}

func (s *ProductStore) FindAll(_ context.Context) ([]domain.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Product, 0, len(s.products))

	for _, p := range s.products {
		result = append(result, p)
	}

	return result, nil
}

func (s *ProductStore) FindByID(_ context.Context, id string) (*domain.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]

	if !ok {
		return nil, domain.ErrProductNotFound
	}

	return &p, nil
}

// seed populates the store with a realistic set of products across VAT categories.
func (s *ProductStore) seed() {
	products := []domain.Product{
		{
			ID:          "prod-001",
			Name:        "Mechanical Keyboard",
			Description: "Full-size mechanical keyboard with Cherry MX switches",
			GrossPrice:  decimal.NewFromFloat(129.99),
			VATRate:     domain.VATStandard,
			Category:    "electronics",
		},
		{
			ID:          "prod-002",
			Name:        "USB-C Hub",
			Description: "7-in-1 USB-C hub with HDMI, SD card reader and PD charging",
			GrossPrice:  decimal.NewFromFloat(49.99),
			VATRate:     domain.VATStandard,
			Category:    "electronics",
		},
		{
			ID:          "prod-003",
			Name:        "Laptop Stand",
			Description: "Aluminium adjustable laptop stand",
			GrossPrice:  decimal.NewFromFloat(39.99),
			VATRate:     domain.VATStandard,
			Category:    "accessories",
		},
		{
			ID:          "prod-004",
			Name:        "Artisan Coffee Beans",
			Description: "Single-origin Ethiopian coffee beans, 500g",
			GrossPrice:  decimal.NewFromFloat(14.90),
			VATRate:     domain.VATReduced,
			Category:    "food",
		},
		{
			ID:          "prod-005",
			Name:        "Whole Grain Bread",
			Description: "Traditional sourdough whole grain loaf",
			GrossPrice:  decimal.NewFromFloat(3.50),
			VATRate:     domain.VATSuper,
			Category:    "food",
		},
	}

	for _, p := range products {
		s.products[p.ID] = p
	}
}
