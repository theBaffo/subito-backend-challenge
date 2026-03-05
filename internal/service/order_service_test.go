package service_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/service"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockProductRepo struct{ mock.Mock }

func (m *mockProductRepo) FindAll(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

type mockOrderRepo struct{ mock.Mock }

func (m *mockOrderRepo) Save(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

var (
	laptopStand = &domain.Product{
		ID:         "prod-001",
		Name:       "Laptop Stand",
		GrossPrice: decimal.NewFromFloat(39.99),
		VATRate:    domain.VATStandard, // 22%
	}

	coffeeBeans = &domain.Product{
		ID:         "prod-002",
		Name:       "Artisan Coffee Beans",
		GrossPrice: decimal.NewFromFloat(14.90),
		VATRate:    domain.VATReduced, // 10%
	}

	bread = &domain.Product{
		ID:         "prod-003",
		Name:       "Whole Grain Bread",
		GrossPrice: decimal.NewFromFloat(3.50),
		VATRate:    domain.VATSuper, // 4%
	}
)

// ---------------------------------------------------------------------------
// CreateOrder tests
// ---------------------------------------------------------------------------

func TestCreateOrder_SingleItem_CalculatesPricesCorrectly(t *testing.T) {
	products := &mockProductRepo{}
	orders := &mockOrderRepo{}
	svc := service.NewOrderService(orders, products)

	products.On("FindByID", mock.Anything, "prod-001").Return(laptopStand, nil)
	orders.On("Save", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-001", Quantity: 2},
		},
	}

	order, err := svc.CreateOrder(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, order)

	assert.Len(t, order.Items, 1)
	item := order.Items[0]

	// Gross total: 39.99 * 2 = 79.98
	assert.True(t, decimal.NewFromFloat(79.98).Equal(item.TotalPrice),
		"expected total price 79.98, got %s", item.TotalPrice)

	// VAT (22%): 79.98 - (79.98 / 1.22) = 79.98 - 65.56 = 14.42
	assert.True(t, decimal.NewFromFloat(14.42).Equal(item.VATAmount),
		"expected VAT 14.42, got %s", item.VATAmount)

	assert.True(t, decimal.NewFromFloat(79.98).Equal(order.TotalPrice))
	assert.True(t, decimal.NewFromFloat(14.42).Equal(order.TotalVAT))
	assert.Equal(t, domain.StatusConfirmed, order.Status)
	assert.NotEmpty(t, order.ID)
}

func TestCreateOrder_MultipleItems_SumsTotalsCorrectly(t *testing.T) {
	products := &mockProductRepo{}
	orders := &mockOrderRepo{}
	svc := service.NewOrderService(orders, products)

	products.On("FindByID", mock.Anything, "prod-001").Return(laptopStand, nil)
	products.On("FindByID", mock.Anything, "prod-002").Return(coffeeBeans, nil)
	orders.On("Save", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-001", Quantity: 1},
			{ProductID: "prod-002", Quantity: 2},
		},
	}

	order, err := svc.CreateOrder(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, order.Items, 2)

	// prod-001: 39.99 * 1 = 39.99
	// prod-002: 14.90 * 2 = 29.80
	// total: 69.79
	expectedTotal := decimal.NewFromFloat(69.79)
	assert.True(t, expectedTotal.Equal(order.TotalPrice),
		"expected total %s, got %s", expectedTotal, order.TotalPrice)
}

func TestCreateOrder_DifferentVATRates_CalculatesEachCorrectly(t *testing.T) {
	products := &mockProductRepo{}
	orders := &mockOrderRepo{}
	svc := service.NewOrderService(orders, products)

	products.On("FindByID", mock.Anything, "prod-003").Return(bread, nil)
	orders.On("Save", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-003", Quantity: 1},
		},
	}

	order, err := svc.CreateOrder(context.Background(), req)
	require.NoError(t, err)

	item := order.Items[0]
	// VAT 4%: 3.50 - (3.50 / 1.04) = 3.50 - 3.37 = 0.13
	assert.True(t, decimal.NewFromFloat(0.13).Equal(item.VATAmount),
		"expected VAT 0.13 for super-reduced rate, got %s", item.VATAmount)
}

func TestCreateOrder_EmptyItems_ReturnsError(t *testing.T) {
	svc := service.NewOrderService(&mockOrderRepo{}, &mockProductRepo{})

	_, err := svc.CreateOrder(context.Background(), service.CreateOrderRequest{Items: []service.OrderItemRequest{}})

	assert.ErrorIs(t, err, domain.ErrEmptyOrder)
}

func TestCreateOrder_ZeroQuantity_ReturnsError(t *testing.T) {
	svc := service.NewOrderService(&mockOrderRepo{}, &mockProductRepo{})

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-001", Quantity: 0},
		},
	}

	_, err := svc.CreateOrder(context.Background(), req)
	assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
}

func TestCreateOrder_NegativeQuantity_ReturnsError(t *testing.T) {
	svc := service.NewOrderService(&mockOrderRepo{}, &mockProductRepo{})

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-001", Quantity: -5},
		},
	}

	_, err := svc.CreateOrder(context.Background(), req)
	assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
}

func TestCreateOrder_ProductNotFound_ReturnsError(t *testing.T) {
	products := &mockProductRepo{}
	svc := service.NewOrderService(&mockOrderRepo{}, products)

	products.On("FindByID", mock.Anything, "prod-999").Return(nil, domain.ErrProductNotFound)

	req := service.CreateOrderRequest{
		Items: []service.OrderItemRequest{
			{ProductID: "prod-999", Quantity: 1},
		},
	}

	_, err := svc.CreateOrder(context.Background(), req)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

// ---------------------------------------------------------------------------
// GetOrder tests
// ---------------------------------------------------------------------------

func TestGetOrder_ExistingOrder_ReturnsOrder(t *testing.T) {
	orders := &mockOrderRepo{}
	svc := service.NewOrderService(orders, &mockProductRepo{})

	expected := &domain.Order{ID: "ord-abc123", Status: domain.StatusConfirmed}
	orders.On("FindByID", mock.Anything, "ord-abc123").Return(expected, nil)

	order, err := svc.GetOrder(context.Background(), "ord-abc123")
	require.NoError(t, err)
	assert.Equal(t, "ord-abc123", order.ID)
}

func TestGetOrder_NonExistentOrder_ReturnsNotFound(t *testing.T) {
	orders := &mockOrderRepo{}
	svc := service.NewOrderService(orders, &mockProductRepo{})

	orders.On("FindByID", mock.Anything, "ord-nope").Return(nil, domain.ErrOrderNotFound)

	_, err := svc.GetOrder(context.Background(), "ord-nope")
	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
}
