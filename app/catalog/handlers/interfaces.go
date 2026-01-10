package handlers

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)
//go:generate mockgen -source=interfaces.go -package=handlers -destination=interfaces_mock.go -mock_names=CatalogService=MockCatalogService
type CatalogService interface {
	GetAllProducts(ctx context.Context, filters *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(ctx context.Context, ProductCode string) (*models.Product, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, Category *models.Category) error
}
