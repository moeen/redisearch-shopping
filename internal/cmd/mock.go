package cmd

import (
	"github.com/moeen/redisearch-shopping/internal/storage/sqlite"
	"github.com/moeen/redisearch-shopping/pkg/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var mockProductsData = []*models.Product{
	{
		Name:  "Bread",
		Price: 4,
	},
	{
		Name:  "Meat",
		Price: 16,
	},
	{
		Name:  "Rice",
		Price: 5,
	},
	{
		Name:  "Eggs",
		Price: 4,
	},
	{
		Name:  "Apples",
		Price: 6,
	},
	{
		Name:  "Potato",
		Price: 4,
	},
	{
		Name:  "Tomato",
		Price: 6,
	},
	{
		Name:  "Onion",
		Price: 4,
	},
	{
		Name:  "Chicken",
		Price: 13,
	},
	{
		Name:  "Milk",
		Price: 1,
	},
}

// mockCommand creates the mock command which populates the database with mock data
func (c *CMD) mockCommand() *cobra.Command {
	mock := &cobra.Command{
		Use:   "mock",
		Long:  "mock will populate the database with mock products data",
		Short: "populate products in database",
		Run:   c.mockRun,
	}

	mock.Flags().StringP("sqlite", "s", "./test.db", "sqlite database file address")

	return mock
}

// mockRun populates the database with mock data
func (c *CMD) mockRun(cmd *cobra.Command, args []string) {
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

	for _, p := range mockProductsData {
		if err := db.AddProduct(p); err != nil {
			c.logger.Error("failed to add product", zap.Error(err))
		}
	}
}
