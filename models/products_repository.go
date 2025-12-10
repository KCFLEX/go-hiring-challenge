package models

import (
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

func (r *ProductsRepository) GetAllProducts(filters *ProductFilter) ([]Product, int64, error) {
	q := r.db.Model(&Product{})

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