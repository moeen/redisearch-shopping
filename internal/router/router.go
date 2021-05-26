package router

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/moeen/redisearch-shopping/graph"
	"github.com/moeen/redisearch-shopping/graph/generated"
	"github.com/moeen/redisearch-shopping/internal/auth"
	"github.com/moeen/redisearch-shopping/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// DefaultPort is used when no port is provided to run the GraphQL server
const DefaultPort = 8080

// setupGraphQLRouter creates the router along with handlers and needed middlewares
func setupGraphQLRouter(mode string, storage storage.Storage, searcher storage.Searcher, logger *zap.Logger) *gin.Engine {
	a := auth.NewAuth(storage)

	gin.SetMode(mode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(ginzap.Ginzap(logger, time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{Storage: storage, Searcher: searcher},
	}))
	router.GET("/", gin.WrapH(playground.Handler("GraphQL playground", "/query")))
	router.POST("/query", a.GinJWTMiddleware, gin.WrapH(srv))

	return router
}

// GraphQLServer creates a http.Server with created GraphQL router
func GraphQLServer(mode string, port int, storage storage.Storage, searcher storage.Searcher, logger *zap.Logger) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupGraphQLRouter(mode, storage, searcher, logger),
	}
}
