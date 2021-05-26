package storage

import "github.com/moeen/redisearch-shopping/pkg/models"

// Storage is the interface used to store all needed data in application
type Storage interface {
	// GetCustomer searches for a customer with an ID and returns it
	GetCustomer(id int) (*models.Customer, error)

	// GetCustomerByEmailAndPassword searches for a customer with an email and returns it
	GetCustomerByEmail(email string) (*models.Customer, error)

	// CreateCustomer creates a new customer with given data
	CreateCustomer(email, name, hash string) (*models.Customer, error)

	// AddToCart adds a product to a customer cart with given quantity
	AddToCart(customerID, productID, quantity int) error

	// RemoveFromCart remove a single product from customer's cart
	RemoveFromCart(customerID, productID int) error

	// GetCartItems returns all items in customer cart
	GetCartItems(customerID int) ([]*models.CartItem, error)

	// AddProduct Will creates the product record in storage
	AddProduct(product *models.Product) error

	// SearchProducts returns all products which has the name in it's name
	// if name is nil, then it returns all the products
	SearchProducts(name *string) ([]*models.Product, error)
}
