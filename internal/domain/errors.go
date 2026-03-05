package domain

import "errors"

// Sentinel errors for the domain layer.
// Handlers inspect these to map them to the correct HTTP status codes.
var (
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	ErrEmptyOrder      = errors.New("order must contain at least one item")
)
