package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/moeen/redisearch-shopping/graph/generated"
	"github.com/moeen/redisearch-shopping/graph/model"
	"github.com/moeen/redisearch-shopping/internal/auth"
	"github.com/moeen/redisearch-shopping/pkg/models"
)

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	c, err := r.Storage.GetCustomerByEmail(input.Email)
	if err != nil {
		return "", errors.New("email or password is wrong")
	}

	if !auth.CheckPasswordHash(input.Password, c.Password) {
		return "", errors.New("email or password is wrong")
	}

	token, err := auth.GenerateToken(int(c.ID), time.Now().Add(auth.DefaultExpirationTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (r *mutationResolver) Register(ctx context.Context, input model.Register) (string, error) {
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return "", err
	}

	c, err := r.Storage.CreateCustomer(input.Email, input.Name, hash)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(int(c.ID), time.Now().Add(auth.DefaultExpirationTime))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (r *mutationResolver) AddToCart(ctx context.Context, input model.AddToCard) (*model.Cart, error) {
	customer, ok := auth.CustomerFromContext(ctx)
	if !ok {
		return nil, errors.New("access denied")
	}

	pID, err := strconv.Atoi(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("inavlid product id: %w", err)
	}

	if err := r.Storage.AddToCart(int(customer.ID), pID, input.Quantity); err != nil {
		return nil, err
	}

	cartItems, err := r.Storage.GetCartItems(int(customer.ID))
	if err != nil {
		return nil, err
	}

	cart := &model.Cart{
		Products: make([]*model.ProductInCart, len(cartItems)),
	}

	for i, ci := range cartItems {
		cart.Products[i] = &model.ProductInCart{
			Product: &model.Product{
				ID:    fmt.Sprintf("%d", ci.ProductID),
				Name:  ci.Product.Name,
				Price: ci.Product.Price,
			},
			Quantity: ci.Quantity,
		}
	}

	return cart, nil
}

func (r *mutationResolver) RemoveFromCart(ctx context.Context, productID string) (*model.Cart, error) {
	customer, ok := auth.CustomerFromContext(ctx)
	if !ok {
		return nil, errors.New("access denied")
	}

	pID, err := strconv.Atoi(productID)
	if err != nil {
		return nil, fmt.Errorf("inavlid product id: %w", err)
	}

	if err := r.Storage.RemoveFromCart(int(customer.ID), pID); err != nil {
		return nil, err
	}

	cartItems, err := r.Storage.GetCartItems(int(customer.ID))
	if err != nil {
		return nil, err
	}

	cart := &model.Cart{
		Products: make([]*model.ProductInCart, len(cartItems)),
	}

	for i, ci := range cartItems {
		cart.Products[i] = &model.ProductInCart{
			Product: &model.Product{
				ID:    fmt.Sprintf("%d", ci.ProductID),
				Name:  ci.Product.Name,
				Price: ci.Product.Price,
			},
			Quantity: ci.Quantity,
		}
	}

	return cart, nil
}

func (r *queryResolver) Products(ctx context.Context, name *string) ([]*model.Product, error) {
	_, ok := auth.CustomerFromContext(ctx)
	if !ok {
		return nil, errors.New("access denied")
	}

	var products []*models.Product
	var err error

	if name == nil || *name == "" {
		products, err = r.Storage.SearchProducts(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get products from storage: %w", err)
		}
	} else {
		products, err = r.Searcher.SearchProducts(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get products from searcher: %w", err)
		}
	}

	res := make([]*model.Product, len(products))
	for i, p := range products {
		res[i] = &model.Product{
			ID:    fmt.Sprintf("%d", p.ID),
			Name:  p.Name,
			Price: p.Price,
		}
	}

	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
