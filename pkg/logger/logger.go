package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewZapLogger will create a new zap logger with given log level
func NewZapLogger(lvl zapcore.Level) *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(config)
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lvl))
}
