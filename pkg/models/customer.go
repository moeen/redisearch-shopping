package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Email    string
	Password string
	Name     string
}
