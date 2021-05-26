package sqlite

import (
	"fmt"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteDatabase is the SQLite implementation of storage.Storage
type SQLiteDatabase struct {
	db *gorm.DB
}

// NewSQLiteDatabase will create a new SQLiteDatabase with given database address
func NewSQLiteDatabase(addr string) (*SQLiteDatabase, error) {
	db, err := gorm.Open(sqlite.Open(addr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &SQLiteDatabase{db}, err
}

// Init will migrate all models needed
func (s *SQLiteDatabase) Init() error {
	err := s.db.AutoMigrate(&models.Customer{}, &models.Product{}, &models.CartItem{})
	if err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	return nil
}

func (s *SQLiteDatabase) GetCustomer(id int) (*models.Customer, error) {
	var c models.Customer
	if err := s.db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, fmt.Errorf("failed to query customer: %w", err)
	}

	return &c, nil
}

func (s *SQLiteDatabase) GetCustomerByEmail(email string) (*models.Customer, error) {
	var c models.Customer
	if err := s.db.Where("email = ?", email).First(&c).Error; err != nil {
		return nil, fmt.Errorf("failed to query customer: %w", err)
	}

	return &c, nil
}

func (s *SQLiteDatabase) CreateCustomer(email, name, hash string) (*models.Customer, error) {
	c := models.Customer{
		Email:    email,
		Password: hash,
		Name:     name,
	}

	if err := s.db.Create(&c).Error; err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &c, nil
}

func (s *SQLiteDatabase) AddToCart(customerID, productID, quantity int) error {
	var cartItem models.CartItem
	err := s.db.Where("customer_id = ? AND product_id = ?", customerID, productID).First(&cartItem).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to query cart item: %w", err)
	}

	if err != gorm.ErrRecordNotFound {
		cartItem.Quantity += quantity
		if err := s.db.Save(&cartItem).Error; err != nil {
			return fmt.Errorf("failed to update cart item: %w", err)
		}
		return nil
	}

	cartItem = models.CartItem{
		CustomerID: customerID,
		Quantity:   quantity,
		ProductID:  productID,
	}

	if err := s.db.Create(&cartItem).Error; err != nil {
		return fmt.Errorf("failed to insert cart item: %w", err)
	}

	return nil
}

func (s *SQLiteDatabase) RemoveFromCart(customerID, productID int) error {
	var cartItem models.CartItem
	if err := s.db.Where("customer_id = ? AND product_id = ?", customerID, productID).First(&cartItem).Error; err != nil {
		return fmt.Errorf("failed to query cart item: %w", err)
	}

	if cartItem.Quantity == 1 {
		if err := s.db.Delete(&cartItem).Error; err != nil {
			return fmt.Errorf("failed to delete cart item: %w", err)
		}

		return nil
	}

	cartItem.Quantity -= 1
	if err := s.db.Save(&cartItem).Error; err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	return nil
}

func (s *SQLiteDatabase) GetCartItems(customerID int) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem

	if err := s.db.Preload("Product").Where("customer_id = ?", customerID).Find(&cartItems).Error; err != nil {
		return nil, fmt.Errorf("failed to query cart items: %w", err)
	}

	return cartItems, nil

}

func (s *SQLiteDatabase) AddProduct(product *models.Product) error {
	if err := s.db.Create(product).Error; err != nil {
		return fmt.Errorf("failed to add product: %w", err)
	}

	return nil
}

func (s *SQLiteDatabase) SearchProducts(name *string) ([]*models.Product, error) {
	var p []*models.Product

	query := s.db
	if name != nil {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *name))
	}

	if err := query.Find(&p).Error; err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}

	return p, nil
}
