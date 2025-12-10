package services

import "github.com/mytheresa/go-hiring-challenge/models"

type productsRepository interface {
	GetAllProducts(filters *models.ProductFilter) ([]models.Product, int64, error)
	GetProductDetailsByCode(ProductCode string) (*models.Product, error)
}
