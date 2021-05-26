package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CustomerID int
	Customer   Customer
	Quantity   int
	ProductID  int
	Product    Product
}
