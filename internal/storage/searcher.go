package storage

import "github.com/moeen/redisearch-shopping/pkg/models"

// Searcher is used to search products
type Searcher interface {
	// SearchProducts returns all products which has the name in it's name
	// if name is nil, then it returns all the products
	SearchProducts(name *string) ([]*models.Product, error)

	// AddProduct will create the given product in searcher and indexes it
	AddProduct(product *models.Product) error
}
