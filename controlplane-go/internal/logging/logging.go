package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func Init() {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder, // e.g., INFO, ERROR
		EncodeTime:  nil,                         // disables timestamp
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	Logger = zap.New(core)
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
