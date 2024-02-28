package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	url := "www.google.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))

	// output
	// {"level":"info","ts":1709085554.4950092,"caller":"log3/main.go:14","msg":"failed to fetch URL","url":"www.google.com","attempt":3,"backoff":1}

	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second)

	// output
	// {"level":"info","ts":1709086291.1179798,"caller":"log3/main.go:23","msg":"failed to fetch URL","url":"www.google.com","attempt":3,"backoff":1}
}
