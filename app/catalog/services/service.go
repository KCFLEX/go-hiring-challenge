package services

import "github.com/mytheresa/go-hiring-challenge/models"

type CatalogService struct {
	productsRepo productsRepository
}

func NewCatalogService(productsRepo productsRepository) *CatalogService {
	return &CatalogService{productsRepo: productsRepo}
}

func (s *CatalogService) GetAllProducts(filters *models.ProductFilter) ([]models.Product, int64, error) {
	products, total, err := s.productsRepo.GetAllProducts(filters)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
