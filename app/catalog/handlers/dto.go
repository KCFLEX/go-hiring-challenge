package handlers

import "github.com/shopspring/decimal"

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductDetailsResponse struct {
	Code     string          `json:"code"`
	Price    decimal.Decimal `json:"price"`
	Category Category        `json:"category"`
	Variant  []Variant       `json:"variants"`
}
