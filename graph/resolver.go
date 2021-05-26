package graph

import "github.com/moeen/redisearch-shopping/internal/storage"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Storage  storage.Storage
	Searcher storage.Searcher
}
