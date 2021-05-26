package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/moeen/redisearch-shopping/internal/storage"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"strings"
)

// JwtContextKey is the key used to store customer in the context
type JwtContextKey struct{}

// Auth is the object used to authenticate incoming requests
type Auth struct {
	storage storage.Storage
}

// NewAuth creates a new Auth with given storage
func NewAuth(storage storage.Storage) *Auth {
	return &Auth{storage: storage}
}

// GinJWTMiddleware is the JWT middleware used for Gin handlers authentication
func (a *Auth) GinJWTMiddleware(ctx *gin.Context) {
	authHeader := strings.Split(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		ctx.Next()
		return
	}

	header := authHeader[1]
	if header == "" {
		ctx.Next()
		return
	}

	customerID, err := ParseToken(header)
	if err != nil {
		ctx.Next()
		return
	}

	customer, err := a.storage.GetCustomer(customerID)
	if err != nil {
		ctx.Next()
		return
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, JwtContextKey{}, customer))
	ctx.Next()
}

// CustomerFromContext searches for the customer in given context
func CustomerFromContext(ctx context.Context) (*models.Customer, bool) {
	c, ok := ctx.Value(JwtContextKey{}).(*models.Customer)
	return c, ok
}
