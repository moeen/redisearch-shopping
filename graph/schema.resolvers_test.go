package graph

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/moeen/redisearch-shopping/graph/model"
	"github.com/moeen/redisearch-shopping/internal/auth"
	"github.com/moeen/redisearch-shopping/internal/storage"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestMutationResolver_Login(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)
	sr := storage.NewMockSearcher(c)

	mr := mutationResolver{&Resolver{
		Storage:  st,
		Searcher: sr,
	}}

	t.Run("test with invalid customer", func(t *testing.T) {
		customer := &models.Customer{
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		st.EXPECT().GetCustomerByEmail(customer.Email).Times(1).Return(nil, errors.New("not found"))

		token, err := mr.Login(context.Background(), model.Login{
			Email:    customer.Email,
			Password: customer.Password,
		})
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("test with wrong password", func(t *testing.T) {
		customer := &models.Customer{
			Email:    "test@test.com",
			Password: "hashed",
			Name:     "test",
		}

		st.EXPECT().GetCustomerByEmail(customer.Email).Times(1).Return(customer, nil)

		token, err := mr.Login(context.Background(), model.Login{
			Email:    customer.Email,
			Password: "test",
		})
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("test successful login", func(t *testing.T) {
		pass := "pass"
		hash, _ := auth.HashPassword(pass)

		customer := &models.Customer{
			Email:    "test@test.com",
			Password: hash,
			Name:     "test",
		}

		st.EXPECT().GetCustomerByEmail(customer.Email).Times(1).Return(customer, nil)

		_, err := mr.Login(context.Background(), model.Login{
			Email:    customer.Email,
			Password: pass,
		})
		assert.NoError(t, err)
	})
}

func TestMutationResolver_Register(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)
	sr := storage.NewMockSearcher(c)

	mr := mutationResolver{&Resolver{
		Storage:  st,
		Searcher: sr,
	}}

	t.Run("test when storage returns an error", func(t *testing.T) {
		input := model.Register{
			Email:    "test@test.com",
			Name:     "test",
			Password: "test",
		}

		st.EXPECT().CreateCustomer(input.Email, input.Name, gomock.Any()).
			Times(1).Return(nil, errors.New("failed"))

		token, err := mr.Register(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("test successful register", func(t *testing.T) {
		input := model.Register{
			Email:    "test@test.com",
			Name:     "test",
			Password: "test",
		}

		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    input.Email,
			Password: input.Password,
			Name:     input.Name,
		}

		st.EXPECT().CreateCustomer(input.Email, input.Name, gomock.Any()).
			Times(1).Return(customer, nil)

		_, err := mr.Register(context.Background(), input)
		assert.NoError(t, err)
	})
}

func TestMutationResolver_AddToCart(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)
	sr := storage.NewMockSearcher(c)

	mr := mutationResolver{&Resolver{
		Storage:  st,
		Searcher: sr,
	}}

	t.Run("test with no customer in ctx", func(t *testing.T) {
		_, err := mr.AddToCart(context.Background(), model.AddToCard{
			ProductID: "1",
			Quantity:  1,
		})

		assert.Error(t, err)
	})

	t.Run("test with invalid product id", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		_, err := mr.AddToCart(ctx, model.AddToCard{
			ProductID: "invalid",
			Quantity:  1,
		})

		assert.Error(t, err)
	})

	t.Run("test when storage.AddToCart returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().AddToCart(int(customer.ID), 1, 1).
			Times(1).Return(errors.New("failed"))

		_, err := mr.AddToCart(ctx, model.AddToCard{
			ProductID: "1",
			Quantity:  1,
		})

		assert.Error(t, err)
	})

	t.Run("test when storage.GetCartItems returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().AddToCart(int(customer.ID), 1, 1).
			Times(1).Return(nil)

		st.EXPECT().GetCartItems(int(customer.ID)).
			Times(1).Return(nil, errors.New("failed"))

		_, err := mr.AddToCart(ctx, model.AddToCard{
			ProductID: "1",
			Quantity:  1,
		})

		assert.Error(t, err)
	})

	t.Run("test successful add to cart", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().AddToCart(int(customer.ID), 1, 1).
			Times(1).Return(nil)

		ret := &model.Cart{
			Products: []*model.ProductInCart{
				{
					Product: &model.Product{
						ID:    "1",
						Name:  "test",
						Price: 1000,
					},
					Quantity: 1,
				},
			},
		}

		var cartItems []*models.CartItem
		for _, p := range ret.Products {
			cartItems = append(cartItems, &models.CartItem{
				CustomerID: 0,
				Customer:   models.Customer{},
				Quantity:   p.Quantity,
				ProductID:  1,
				Product: models.Product{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  p.Product.Name,
					Price: p.Product.Price,
				},
			})
		}

		st.EXPECT().GetCartItems(int(customer.ID)).
			Times(1).Return(cartItems, nil)

		items, err := mr.AddToCart(ctx, model.AddToCard{
			ProductID: "1",
			Quantity:  1,
		})
		assert.NoError(t, err)
		assert.Equal(t, ret, items)
	})
}

func TestMutationResolver_RemoveFromCart(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)
	sr := storage.NewMockSearcher(c)

	mr := mutationResolver{&Resolver{
		Storage:  st,
		Searcher: sr,
	}}

	t.Run("test with no customer in ctx", func(t *testing.T) {
		_, err := mr.RemoveFromCart(context.Background(), "1")

		assert.Error(t, err)
	})

	t.Run("test with invalid product id", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		_, err := mr.RemoveFromCart(ctx, "invalid")

		assert.Error(t, err)
	})

	t.Run("test when storage.RemoveFromCart returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().RemoveFromCart(int(customer.ID), 1).
			Times(1).Return(errors.New("failed"))

		_, err := mr.RemoveFromCart(ctx, "1")

		assert.Error(t, err)
	})

	t.Run("test when storage.GetCartItems returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().RemoveFromCart(int(customer.ID), 1).
			Times(1).Return(nil)

		st.EXPECT().GetCartItems(int(customer.ID)).
			Times(1).Return(nil, errors.New("failed"))

		_, err := mr.RemoveFromCart(ctx, "1")

		assert.Error(t, err)
	})

	t.Run("test successful remove from cart", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().RemoveFromCart(int(customer.ID), 1).
			Times(1).Return(nil)

		ret := &model.Cart{
			Products: []*model.ProductInCart{
				{
					Product: &model.Product{
						ID:    "1",
						Name:  "test",
						Price: 1000,
					},
					Quantity: 1,
				},
			},
		}

		var cartItems []*models.CartItem
		for _, p := range ret.Products {
			cartItems = append(cartItems, &models.CartItem{
				CustomerID: 0,
				Customer:   models.Customer{},
				Quantity:   p.Quantity,
				ProductID:  1,
				Product: models.Product{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  p.Product.Name,
					Price: p.Product.Price,
				},
			})
		}

		st.EXPECT().GetCartItems(int(customer.ID)).
			Times(1).Return(cartItems, nil)

		items, err := mr.RemoveFromCart(ctx, "1")
		assert.NoError(t, err)
		assert.Equal(t, ret, items)
	})
}

func TestQueryResolver_Products(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)
	sr := storage.NewMockSearcher(c)

	r := queryResolver{&Resolver{
		Storage:  st,
		Searcher: sr,
	}}

	t.Run("test with no customer in ctx", func(t *testing.T) {
		name := "product"
		_, err := r.Products(context.Background(), &name)

		assert.Error(t, err)
	})

	t.Run("test when storage.SearchProducts returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		st.EXPECT().SearchProducts(nil).Times(1).Return(nil, errors.New("failed"))

		r, err := r.Products(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("test successful search when name is null", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		products := []*models.Product{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "test1",
				Price: 10,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Name:  "test2",
				Price: 20,
			},
		}

		st.EXPECT().SearchProducts(nil).Times(1).Return(products, nil)

		r, err := r.Products(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(r))
	})

	t.Run("test successful search when name length is zero", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		products := []*models.Product{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "test1",
				Price: 10,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Name:  "test2",
				Price: 20,
			},
		}

		name := ""

		st.EXPECT().SearchProducts(&name).Times(1).Return(products, nil)

		r, err := r.Products(ctx, &name)
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(r))
	})

	t.Run("test when searcher.SearchProducts returns an error", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		name := "test"
		sr.EXPECT().SearchProducts(&name).Times(1).Return(nil, errors.New("failed"))

		r, err := r.Products(ctx, &name)
		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("test successful search when name is provided", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), auth.JwtContextKey{}, customer)

		products := []*models.Product{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "test1",
				Price: 10,
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Name:  "test2",
				Price: 20,
			},
		}

		name := "test"

		sr.EXPECT().SearchProducts(&name).Times(1).Return(products, nil)

		r, err := r.Products(ctx, &name)
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(r))
	})
}
