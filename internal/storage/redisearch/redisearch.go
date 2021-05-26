package redisearch

import (
	"fmt"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/moeen/redisearch-shopping/internal/storage"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"gorm.io/gorm"
	"strconv"
)

// RediSearch is the RediSearch implementation of storage.Searcher
type RediSearch struct {
	rs      *redisearch.Client
	storage storage.Storage
}

// NewRediSearch will create a new RediSearch with given required params
func NewRediSearch(address, index string, s storage.Storage) *RediSearch {
	c := redisearch.NewClient(address, index)
	return &RediSearch{rs: c, storage: s}
}

// Init will create the schema and adds all products to RediSearch
func (r *RediSearch) Init() error {
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewNumericFieldOptions("id", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("name", redisearch.TextFieldOptions{Sortable: true, NoIndex: false})).
		AddField(redisearch.NewNumericField("price"))

	r.rs.Drop()

	if err := r.rs.CreateIndex(sc); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	products, err := r.storage.SearchProducts(nil)
	if err != nil {
		return fmt.Errorf("failed to get products from storage: %s", err)
	}

	for _, p := range products {
		if err := r.AddProduct(p); err != nil {
			return fmt.Errorf("failed to add product to searcher: %w", err)
		}
	}

	return nil
}

func (r *RediSearch) SearchProducts(name *string) ([]*models.Product, error) {
	docs, total, err := r.rs.Search(redisearch.NewQuery(fmt.Sprintf("@name:%s*", *name)).
		SetReturnFields("id", "name", "price"))

	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if total == 0 {
		return []*models.Product{}, nil
	}

	res := make([]*models.Product, total)
	for i, p := range docs {
		id, err := strconv.Atoi(p.Properties["id"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to convert id to str")
		}

		price, err := strconv.Atoi(p.Properties["price"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to convert price to str")
		}

		res[i] = &models.Product{
			Model: gorm.Model{
				ID: uint(id),
			},
			Name:  p.Properties["name"].(string),
			Price: price,
		}
	}

	return res, nil
}

func (r *RediSearch) AddProduct(product *models.Product) error {
	doc := redisearch.NewDocument(fmt.Sprintf("product:%d", product.ID), 1.0)
	doc.Set("id", product.ID).
		Set("name", product.Name).
		Set("price", product.Price)

	if err := r.rs.Index(doc); err != nil {
		return fmt.Errorf("failed to create doc: %w", err)
	}

	return nil
}
