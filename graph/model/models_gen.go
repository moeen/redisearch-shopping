// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddToCard struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Cart struct {
	Products []*ProductInCart `json:"products"`
}

type Customer struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Cart     *Cart  `json:"cart"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductInCart struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type Register struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
