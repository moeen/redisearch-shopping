package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/moeen/redisearch-shopping/internal/storage"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAuth(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)

	auth := NewAuth(st)
	assert.Same(t, st, auth.storage)
}

func TestAuth_GinJWTMiddleware(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	st := storage.NewMockStorage(c)

	auth := NewAuth(st)
	router := gin.New()
	router.GET("/test", auth.GinJWTMiddleware, func(ctx *gin.Context) {
		customer, ok := CustomerFromContext(ctx.Request.Context())
		if !ok {
			ctx.String(http.StatusForbidden, "access denied")
			return
		}
		ctx.String(http.StatusOK, customer.Email)
	})

	t.Run("test with no auth header", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("test with no jwt token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", "Bearer ")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("test with invalid jwt token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", "Bearer test")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("test with invalid customer id", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 10,
			},
			Email: "test@test.com",
		}

		st.EXPECT().GetCustomer(int(customer.ID)).Times(1).Return(nil, errors.New("not found"))

		token, err := GenerateToken(int(customer.ID), time.Now().Add(time.Hour))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("test with valid jwt token", func(t *testing.T) {
		customer := &models.Customer{
			Model: gorm.Model{
				ID: 10,
			},
			Email: "test@test.com",
		}

		st.EXPECT().GetCustomer(int(customer.ID)).Times(1).Return(customer, nil)

		token, err := GenerateToken(int(customer.ID), time.Now().Add(time.Hour))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		body, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, customer.Email, string(body))
	})
}

func TestCustomerFromContext(t *testing.T) {
	t.Run("test with customer in ctx", func(t *testing.T) {
		c := &models.Customer{
			Email:    "test@test.com",
			Password: "test",
			Name:     "test",
		}

		ctx := context.WithValue(context.Background(), JwtContextKey{}, c)

		rc, ok := CustomerFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, c, rc)
	})

	t.Run("test with no customer in ctx", func(t *testing.T) {
		ctx := context.Background()

		rc, ok := CustomerFromContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, rc)
	})
}
