package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents a single product line within an order,
// capturing a snapshot of pricing at the time of purchase.
type OrderItem struct {
	ProductID  string          `json:"product_id"`
	Name       string          `json:"name"`
	Quantity   int             `json:"quantity"`
	UnitPrice  decimal.Decimal `json:"unit_price"`  // Gross price per unit
	TotalPrice decimal.Decimal `json:"total_price"` // Gross price * quantity
	VATRate    decimal.Decimal `json:"vat_rate"`    // e.g. 0.22
	VATAmount  decimal.Decimal `json:"vat_amount"`  // Total VAT for this line
}

// Order is the aggregate root for a purchase transaction.
type Order struct {
	ID         string          `json:"id"`
	Items      []OrderItem     `json:"items"`
	TotalPrice decimal.Decimal `json:"total_price"` // Sum of all item gross prices
	TotalVAT   decimal.Decimal `json:"total_vat"`   // Sum of all item VAT amounts
	Status     OrderStatus     `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
}
