package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootCommand returns the base command of the application
func (c *CMD) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "redisearch-shopping",
		Long:  "redisearch-shopping is a shopping app with Redisearch",
		Short: "redisearch-shopping app",
		Run:   c.rootRun,
	}
}

// rootRun only prints a message
func (*CMD) rootRun(cmd *cobra.Command, args []string) {
	fmt.Println("run `redisearch-shopping serve` to run the http server")
}
