package domain

import "github.com/shopspring/decimal"

// VATRate represents a named VAT category used in Italy.
// Prices are stored gross (tax-inclusive), as per EU convention.
type VATRate struct {
	Name string
	Rate decimal.Decimal // e.g. 0.22 for 22%
}

// Common Italian VAT rates.
var (
	VATStandard = VATRate{Name: "standard", Rate: decimal.NewFromFloat(0.22)}      // Most goods
	VATReduced  = VATRate{Name: "reduced", Rate: decimal.NewFromFloat(0.10)}       // Food, tourism
	VATSuper    = VATRate{Name: "super_reduced", Rate: decimal.NewFromFloat(0.04)} // Essential goods
)

// Product is the core entity representing a purchasable item.
type Product struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	GrossPrice  decimal.Decimal `json:"gross_price"` // Tax-inclusive price
	VATRate     VATRate         `json:"vat_rate"`
	Category    string          `json:"category"`
}

// NetPrice returns the price excluding VAT.
// Formula: gross / (1 + rate)
func (p Product) NetPrice() decimal.Decimal {
	divisor := decimal.NewFromInt(1).Add(p.VATRate.Rate)
	return p.GrossPrice.Div(divisor).RoundBank(2)
}

// VATAmount returns the VAT component of the gross price.
func (p Product) VATAmount() decimal.Decimal {
	return p.GrossPrice.Sub(p.NetPrice()).RoundBank(2)
}
