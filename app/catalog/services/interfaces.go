package services

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type productsRepository interface {
	GetAllProducts(ctx context.Context, filters *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetailsByCode(ctx context.Context, ProductCode string) (*models.Product, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	CreateCategory(ctx context.Context, Category *models.Category) error
}
