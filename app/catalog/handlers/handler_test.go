package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/services"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestHandleGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatService := NewMockCatalogService(ctrl)

	type fields struct {
		srv CatalogService
	}

	tests := []struct {
		name       string
		fields     fields
		url        string
		setup      func()
		wantStatus int
	}{
		{
			name: "suceessfull products retrival",
			fields: fields{
				srv: mockCatService,
			},
			url: "/catalog?offset=0&limit=10",
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllProducts(
						gomock.Any(),
						gomock.AssignableToTypeOf(&models.ProductFilter{}),
					).
					Return(
						[]models.Product{
							{
								Code:  "P1",
								Price: decimal.NewFromFloat(100),
								Category: models.Category{
									Code: "C1",
									Name: "Shoes",
								},
							},
						},
						int64(1),
						nil,
					)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "internal service error",
			fields: fields{
				srv: mockCatService,
			},
			url: "/catalog",
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllProducts(
						gomock.Any(),
						gomock.AssignableToTypeOf(&models.ProductFilter{}),
					).
					Return(nil, int64(0), errors.New("db down"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "category filter applied",
			fields: fields{
				srv: mockCatService,
			},
			url: "/catalog?category=shoes",
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllProducts(
						gomock.Any(),
						gomock.AssignableToTypeOf(&models.ProductFilter{}),
					).
					Return([]models.Product{}, int64(0), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "limit greater than max is capped",
			fields: fields{
				srv: mockCatService,
			},
			url: "/catalog?limit=1000",
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllProducts(
						gomock.Any(),
						gomock.AssignableToTypeOf(&models.ProductFilter{}),
					).
					Return([]models.Product{}, int64(0), nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "price filter applied",
			fields: fields{
				srv: mockCatService,
			},
			url: "/catalog?price_lt=50.5",
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllProducts(
						gomock.Any(),
						gomock.AssignableToTypeOf(&models.ProductFilter{}),
					).
					Return([]models.Product{}, int64(0), nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			hh := &CatalogHandler{service: tt.fields.srv}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			rec := httptest.NewRecorder()

			hh.HandleGet(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}

}

func TestHandleGetProdDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatService := NewMockCatalogService(ctrl)

	type fields struct {
		srv CatalogService
	}

	tests := []struct {
		name       string
		fields     fields
		code       string
		setup      func()
		wantStatus int
	}{
		{
			name: "successful product retrieval",
			fields: fields{
				srv: mockCatService,
			},
			code: "P1",
			setup: func() {
				mockCatService.
					EXPECT().
					GetProductDetails(gomock.Any(), "P1").
					Return(
						&models.Product{
							Code:  "P1",
							Price: decimal.NewFromFloat(100),
							Category: models.Category{
								Code: "C1",
								Name: "Shoes",
							},
							Variants: []models.Variant{
								{
									Name:  "Red",
									SKU:   "SKU-RED",
									Price: decimal.NewFromFloat(90),
								},
							},
						},
						nil,
					)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing product code",
			fields: fields{
				srv: mockCatService,
			},
			code:       "",
			setup:      func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "product not found",
			fields: fields{
				srv: mockCatService,
			},
			code: "P404",
			setup: func() {
				mockCatService.
					EXPECT().
					GetProductDetails(gomock.Any(), "P404").
					Return(nil, services.ErrProductNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal service error",
			fields: fields{
				srv: mockCatService,
			},
			code: "P500",
			setup: func() {
				mockCatService.
					EXPECT().
					GetProductDetails(gomock.Any(), "P500").
					Return(nil, errors.New("db failure"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			h := &CatalogHandler{
				service: tt.fields.srv,
			}

			req := httptest.NewRequest(
				http.MethodGet,
				"/catalog/"+tt.code,
				nil,
			)

			if tt.code != "" {
				req.SetPathValue("code", tt.code)
			}

			rec := httptest.NewRecorder()

			h.HandleGetProdDetails(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandleGetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatService := NewMockCatalogService(ctrl)

	type fields struct {
		srv CatalogService
	}

	tests := []struct {
		name       string
		fields     fields
		setup      func()
		wantStatus int
	}{
		{
			name: "successful categories retrieval",
			fields: fields{
				srv: mockCatService,
			},
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(
						[]models.Category{
							{
								Code: "C1",
								Name: "Shoes",
							},
							{
								Code: "C2",
								Name: "Jackets",
							},
						},
						nil,
					)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "nil categories slice returned",
			fields: fields{
				srv: mockCatService,
			},
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service error",
			fields: fields{
				srv: mockCatService,
			},
			setup: func() {
				mockCatService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			h := &CatalogHandler{
				service: tt.fields.srv,
			}

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			rec := httptest.NewRecorder()

			h.HandleGetAllCategories(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestHandleCreateNewCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCatService := NewMockCatalogService(ctrl)

	type fields struct {
		srv CatalogService
	}

	tests := []struct {
		name       string
		fields     fields
		body       string
		setup      func()
		wantStatus int
	}{
		{
			name: "successful category creation",
			fields: fields{
				srv: mockCatService,
			},
			body: `{"code":"C1","name":"Shoes"}`,
			setup: func() {
				mockCatService.
					EXPECT().
					CreateCategory(
						gomock.Any(),
						&models.Category{
							Code: "C1",
							Name: "Shoes",
						},
					).
					Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid json body",
			fields: fields{
				srv: mockCatService,
			},
			body:       `{"code":`,
			setup:      nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing category code",
			fields: fields{
				srv: mockCatService,
			},
			body:       `{"name":"Shoes"}`,
			setup:      nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing category name",
			fields: fields{
				srv: mockCatService,
			},
			body:       `{"code":"C1"}`,
			setup:      nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			fields: fields{
				srv: mockCatService,
			},
			body: `{"code":"C1","name":"Shoes"}`,
			setup: func() {
				mockCatService.
					EXPECT().
					CreateCategory(
						gomock.Any(),
						&models.Category{
							Code: "C1",
							Name: "Shoes",
						},
					).
					Return(errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			h := &CatalogHandler{
				service: tt.fields.srv,
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/categories",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.HandleCreateNewCategory(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
