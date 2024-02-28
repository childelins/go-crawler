package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 定制化logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Hello from zap logger")

	// output
	// {"level":"info","timestamp":"2024-02-28T10:17:39+08:00","caller":"log4/main.go:23","msg":"Hello from zap logger"}
}
