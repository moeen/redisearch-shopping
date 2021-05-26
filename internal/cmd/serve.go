package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/moeen/redisearch-shopping/internal/router"
	"github.com/moeen/redisearch-shopping/internal/storage/redisearch"
	"github.com/moeen/redisearch-shopping/internal/storage/sqlite"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// serveCommand creates the serve command which runs the GraphQL server
func (c *CMD) serveCommand() *cobra.Command {
	serve := &cobra.Command{
		Use:   "serve",
		Long:  "serve runs the http server along with GraphQL",
		Short: "run server",
		Run:   c.serveRun,
	}
	serve.Flags().IntP("port", "p", router.DefaultPort, "http server port")
	serve.Flags().StringP("mode", "m", gin.DebugMode, "router mode")
	serve.Flags().StringP("sqlite", "s", "./test.db", "sqlite database file address")

	return serve
}

// serveRun parses all needed flags, creates all the dependencies and runs the GraphQL server
func (c *CMD) serveRun(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		c.logger.Fatal("failed to get the port", zap.Error(err))
	}

	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		c.logger.Fatal("failed to get the mode", zap.Error(err))
	}

	addr, err := cmd.Flags().GetString("sqlite")
	if err != nil {
		c.logger.Fatal("failed to get the sqlite address", zap.Error(err))
	}

	db, err := sqlite.NewSQLiteDatabase(addr)
	if err != nil {
		c.logger.Fatal("failed to create sqlite db", zap.Error(err))
	}
	if err := db.Init(); err != nil {
		c.logger.Fatal("failed to init database", zap.Error(err))
	}

	rs := redisearch.NewRediSearch("127.0.0.1:6379", "test", db)
	if err := rs.Init(); err != nil {
		c.logger.Fatal("failed to init RediSearch", zap.Error(err))
	}

	restServer := router.GraphQLServer(mode, port, db, rs, c.logger.Named("router"))
	c.logger.Fatal(restServer.ListenAndServe().Error())
}
