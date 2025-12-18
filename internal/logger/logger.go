package logger

import "go.uber.org/zap"

// NewLogger creates a production-ready Zap logger.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}


