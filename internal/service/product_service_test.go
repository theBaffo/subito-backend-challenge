package service_test

import (
	"context"
	"testing"

	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Note: mockProductRepo and mockOrderRepo are already declared in order_service_test.go.
// They are reused here within the same package.

var catalogueFixture = []domain.Product{
	{
		ID:         "prod-001",
		Name:       "Laptop Stand",
		GrossPrice: decimal.NewFromFloat(39.99),
		VATRate:    domain.VATStandard,
		Category:   "accessories",
	},
	{
		ID:         "prod-002",
		Name:       "Artisan Coffee Beans",
		GrossPrice: decimal.NewFromFloat(14.90),
		VATRate:    domain.VATReduced,
		Category:   "food",
	},
}

func TestGetAllProducts_ReturnsCatalogue(t *testing.T) {
	repo := &mockProductRepo{}
	svc := service.NewProductService(repo)

	repo.On("FindAll", mock.Anything).Return(catalogueFixture, nil)

	products, err := svc.GetAllProducts(context.Background())
	require.NoError(t, err)
	assert.Len(t, products, 2)
	repo.AssertExpectations(t)
}

func TestGetAllProducts_EmptyCatalogue_ReturnsEmptySlice(t *testing.T) {
	repo := &mockProductRepo{}
	svc := service.NewProductService(repo)

	repo.On("FindAll", mock.Anything).Return([]domain.Product{}, nil)

	products, err := svc.GetAllProducts(context.Background())
	require.NoError(t, err)
	assert.Empty(t, products)
}

func TestGetProduct_ExistingID_ReturnsProduct(t *testing.T) {
	repo := &mockProductRepo{}
	svc := service.NewProductService(repo)

	repo.On("FindByID", mock.Anything, "prod-001").Return(&catalogueFixture[0], nil)

	product, err := svc.GetProduct(context.Background(), "prod-001")
	require.NoError(t, err)
	assert.Equal(t, "prod-001", product.ID)
	assert.Equal(t, "Laptop Stand", product.Name)
	repo.AssertExpectations(t)
}

func TestGetProduct_NonExistentID_ReturnsNotFound(t *testing.T) {
	repo := &mockProductRepo{}
	svc := service.NewProductService(repo)

	repo.On("FindByID", mock.Anything, "prod-999").Return(nil, domain.ErrProductNotFound)

	_, err := svc.GetProduct(context.Background(), "prod-999")
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}
