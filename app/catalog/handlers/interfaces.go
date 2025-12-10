package handlers

import "github.com/mytheresa/go-hiring-challenge/models"

type CatalogService interface {
	GetAllProducts(filters *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetails(ProductCode string) (*models.Product, error)
	GetAllCategories() ([]models.Category, error)
	CreateCategory(Category *models.Category) error
}
