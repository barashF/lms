package logger

import (
	"go.uber.org/zap"
)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter() (*ZapAdapter, error) {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Encoding = "console"

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &ZapAdapter{logger: zapLogger}, nil
}

func (z *ZapAdapter) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Sync() error {
	return z.logger.Sync()
}

func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
