package services

import "github.com/mytheresa/go-hiring-challenge/models"

type productsRepository interface {
	GetAllProducts() ([]models.Product, error)
}
