package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// 定制化logger
	w := &lumberjack.Logger{
		Filename:   "./my.log",
		MaxSize:    500, // 日志的最大大小，以M为单位
		MaxBackups: 3,   // 保留的旧日志文件的最大数量
		MaxAge:     28,  // 保留旧日志文件的最大天数
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(zapcore.AddSync(w)),
		zap.InfoLevel,
	)

	logger := zap.New(core)
	defer logger.Sync()

	logger.Info("Hello from zap logger with lumberjack")
}
