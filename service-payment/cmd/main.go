package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Encoding = "console"

	zapLogger, err := config.Build()
	zapLogger.Info("test")
	zapLogger.Info("password")
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
}
