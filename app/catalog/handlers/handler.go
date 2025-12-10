package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

type CatalogHandler struct {
	service CatalogService
}

func NewCatalogHandler(service CatalogService) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {

	offset, err := getQueryInt(r, "offset", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getQueryInt(r, "limit", 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if limit > 100 {
		limit = 100
	}

	if limit < 1 {
		limit = 1
	}

	category := r.URL.Query().Get("category")
	priceStr := r.URL.Query().Get("price_lt")

	var priceLt *decimal.Decimal
	if priceStr != "" {
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("Error parsing price: ", err)
			return
		}
		priceLt = &price
	}

	filters := &models.ProductFilter{
		Offset:       offset,
		Limit:        limit,
		CategoryCode: &category,
		PriceLt: priceLt,
	}

	res, prodTotal, err := h.service.GetAllProducts(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
		Total:    prodTotal,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
