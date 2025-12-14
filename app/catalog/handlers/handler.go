package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
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
	ctx := r.Context()
	offset, err := getQueryInt(r, "offset", 0)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, err := getQueryInt(r, "limit", 10)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
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
			api.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		priceLt = &price
	}

	filters := &models.ProductFilter{
		Offset:       offset,
		Limit:        limit,
		CategoryCode: &category,
		PriceLt:      priceLt,
	}

	res, prodTotal, err := h.service.GetAllProducts(ctx, filters)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
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

	response := Response{
		Products: products,
		Total:    prodTotal,
	}
	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetProdDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	codeParam := r.PathValue("code")

	product, err := h.service.GetProductDetails(ctx, codeParam)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	responseDTO := ProductDetailsResponse{
		Code:  product.Code,
		Price: product.Price,
		Category: Category{
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
		Variant: make([]Variant, len(product.Variants)),
	}

	for i, v := range product.Variants {
		responseDTO.Variant[i] = Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		}
	}

	api.OKResponse(w, responseDTO)

}

func (h *CatalogHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.service.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]Category, len(categories))
	for i, c := range categories {
		response[i] = Category{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleCreateNewCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newCategory := &models.Category{
		Code: category.Code,
		Name: category.Name,
	}

	if err := h.service.CreateCategory(ctx, newCategory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
