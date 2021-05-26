package main

import (
	"github.com/moeen/redisearch-shopping/internal/cmd"
	"github.com/moeen/redisearch-shopping/pkg/logger"
	"go.uber.org/zap/zapcore"
)

// default application log level
const logLevel = zapcore.DebugLevel

func main() {
	l := logger.NewZapLogger(logLevel).Named("main")
	c := cmd.NewCMD(l.Named("cmd"))

	if err := c.Execute(); err != nil {
		panic(err)
	}
}
