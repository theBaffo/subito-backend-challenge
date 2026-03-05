package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/theBaffo/subito-backend-challenge/internal/domain"
	"github.com/theBaffo/subito-backend-challenge/internal/repository"
)

// CreateOrderRequest is the input DTO for creating a new order.
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

// OrderItemRequest specifies a product and how many units to purchase.
type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// OrderService contains the business logic for order operations.
// It depends only on repository interfaces, keeping it storage-agnostic.
type OrderService struct {
	orders   repository.OrderRepository
	products repository.ProductRepository
}

// NewOrderService constructs an OrderService with its required dependencies.
func NewOrderService(orders repository.OrderRepository, products repository.ProductRepository) *OrderService {
	return &OrderService{
		orders:   orders,
		products: products,
	}
}

// CreateOrder validates the request, prices each item, computes totals,
// persists the order, and returns the fully populated Order.
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	items, totalPrice, totalVAT, err := s.buildOrderItems(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:         "ord-" + uuid.New().String()[:8],
		Items:      items,
		TotalPrice: totalPrice,
		TotalVAT:   totalVAT,
		Status:     domain.StatusConfirmed,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.orders.Save(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by its ID.
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return s.orders.FindByID(ctx, id)
}

// buildOrderItems resolves each requested item against the product catalogue,
// calculates per-line pricing, and accumulates the order totals.
func (s *OrderService) buildOrderItems(
	ctx context.Context,
	requests []OrderItemRequest,
) ([]domain.OrderItem, decimal.Decimal, decimal.Decimal, error) {

	items := make([]domain.OrderItem, 0, len(requests))
	totalPrice := decimal.Zero
	totalVAT := decimal.Zero

	for _, req := range requests {
		product, err := s.products.FindByID(ctx, req.ProductID)
		if err != nil {
			return nil, decimal.Zero, decimal.Zero, err
		}

		qty := decimal.NewFromInt(int64(req.Quantity))
		lineGross := product.GrossPrice.Mul(qty).RoundBank(2)

		// VAT component: gross - (gross / (1 + rate))
		divisor := decimal.NewFromInt(1).Add(product.VATRate.Rate)
		lineNet := lineGross.Div(divisor).RoundBank(2)
		lineVAT := lineGross.Sub(lineNet).RoundBank(2)

		items = append(items, domain.OrderItem{
			ProductID:  product.ID,
			Name:       product.Name,
			Quantity:   req.Quantity,
			UnitPrice:  product.GrossPrice,
			TotalPrice: lineGross,
			VATRate:    product.VATRate.Rate,
			VATAmount:  lineVAT,
		})

		totalPrice = totalPrice.Add(lineGross)
		totalVAT = totalVAT.Add(lineVAT)
	}

	return items, totalPrice.RoundBank(2), totalVAT.RoundBank(2), nil
}

// validateRequest performs input validation before any business logic runs.
func validateRequest(req CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return domain.ErrEmptyOrder
	}
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return domain.ErrInvalidQuantity
		}
	}
	return nil
}
