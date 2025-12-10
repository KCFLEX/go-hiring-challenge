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

func (s *CatalogService) GetProductDetails(ProductCode string) (*models.Product, error) {
	product, err := s.productsRepo.GetProductDetailsByCode(ProductCode)
	if err != nil {
		return &models.Product{}, err
	}

	return product, nil
}

func (s *CatalogService) GetAllCategories() ([]models.Category, error) {
	categories, err := s.productsRepo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CatalogService) CreateCategory(Category *models.Category) error {
	err := s.productsRepo.CreateCategory(Category)
	if err != nil {
		return err
	}
	return nil
}
