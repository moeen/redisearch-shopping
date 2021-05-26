package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// CMD is the struct responsible for all commands in the app
type CMD struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

// NewCMD receives a logger and creates the CMD along with all app commands
func NewCMD(logger *zap.Logger) *CMD {
	c := &CMD{
		logger: logger,
	}

	root := c.rootCommand()
	serve := c.serveCommand()
	mock := c.mockCommand()

	root.AddCommand(serve)
	root.AddCommand(mock)

	c.cmd = root

	return c
}

// Execute runs the CMD
func (c *CMD) Execute() error {
	return c.cmd.Execute()
}
