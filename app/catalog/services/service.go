package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")
var ErrNoProducts = errors.New("no products found")

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

	if len(products) == 0 {
		return nil, 0, ErrNoProducts
	}

	return products, total, nil
}

func (s *CatalogService) GetProductDetails(ctx context.Context, ProductCode string) (*models.Product, error) {
	product, err := s.productsRepo.GetProductDetailsByCode(ctx, ProductCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product %s not found: %w", ProductCode, ErrProductNotFound)
		}
		return nil, fmt.Errorf("failed to get product %s: %w", ProductCode, err)
	}

	return product, nil
}

func (s *CatalogService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := s.productsRepo.GetAllCategories(ctx)
	if err != nil {
		 return nil, fmt.Errorf("get all categories failed: %w", err)
	}

	return categories, nil
}

func (s *CatalogService) CreateCategory(ctx context.Context, Category *models.Category) error {
	err := s.productsRepo.CreateCategory(ctx, Category)
	if err != nil {
		return fmt.Errorf("create category failed: %w", err)
	}
	return nil
}
