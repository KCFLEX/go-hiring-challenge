package services

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogService struct {
	productsRepo productsRepository
}

func NewCatalogService(productsRepo productsRepository) *CatalogService {
	return &CatalogService{productsRepo: productsRepo}
}

func (s *CatalogService) GetAllProducts(ctx context.Context, filters *models.ProductFilter) ([]models.Product, int64, error) {
	products, total, err := s.productsRepo.GetAllProducts(ctx, filters)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *CatalogService) GetProductDetails(ctx context.Context, ProductCode string) (*models.Product, error) {
	product, err := s.productsRepo.GetProductDetailsByCode(ctx, ProductCode)
	if err != nil {
		return &models.Product{}, err
	}

	return product, nil
}

func (s *CatalogService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := s.productsRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CatalogService) CreateCategory(ctx context.Context, Category *models.Category) error {
	err := s.productsRepo.CreateCategory(ctx, Category)
	if err != nil {
		return err
	}
	return nil
}
