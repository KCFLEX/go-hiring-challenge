package services

import "github.com/mytheresa/go-hiring-challenge/models"

type CatalogService struct {
	productsRepo productsRepository
}

func NewCatalogService(productsRepo productsRepository) *CatalogService {
	return &CatalogService{productsRepo: productsRepo}
}

func (s *CatalogService) GetAllProducts() ([]models.Product, error) {
	return s.productsRepo.GetAllProducts()
}
