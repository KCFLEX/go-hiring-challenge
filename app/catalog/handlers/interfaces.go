package handlers

import "github.com/mytheresa/go-hiring-challenge/models"

type CatalogService interface {
	GetAllProducts() ([]models.Product, error)
}
