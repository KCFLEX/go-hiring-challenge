package models

import "github.com/shopspring/decimal"

type ProductFilter struct {
	Offset   int
	Limit    int
	CategoryCode *string
	PriceLt  *decimal.Decimal
}
