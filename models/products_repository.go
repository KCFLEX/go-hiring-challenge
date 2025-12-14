package models

import (
	"context"

	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(ctx context.Context, filters *ProductFilter) ([]Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&Product{})

	if filters.CategoryCode != nil && *filters.CategoryCode != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", *filters.CategoryCode)
	}

	if filters.PriceLt != nil {
		q = q.Where("products.price < ?", *filters.PriceLt)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []Product
	if err := q.Preload("Variants").Preload("Category").
		Offset(filters.Offset).Limit(filters.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductsRepository) GetProductDetailsByCode(ctx context.Context, ProductCode string) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).Preload("Variants").
		Preload("Category").
		Where("code = ?", ProductCode).First(&product).Error

	if err != nil {
		return nil, err
	}

	for i := range product.Variants {
		if product.Variants[i].Price.IsZero() {
			product.Variants[i].Price = product.Price
		}
	}

	return &product, nil
}

func (r *ProductsRepository) GetAllCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsRepository) CreateCategory(ctx context.Context, Category *Category) error {
	if err := r.db.WithContext(ctx).Create(&Category).Error; err != nil {
		return err
	}
	return nil
}
