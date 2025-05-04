package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Initialize sets up the global logger instance with the appropriate configuration
// based on the provided environment. It uses a production configuration by default,
// but switches to a development configuration with colorized log levels if the
// environment is not "production". The logger includes a predefined "service" field
// with the value "links-service-write".
//
// Parameters:
//   - environment: A string indicating the current environment (e.g., "production",
//     "development"). Determines the logging configuration to use.
func Initialize(environment string) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if environment != "production" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	Log, _ = config.Build(
		zap.Fields(zap.String("service", "links-service-write")),
	)
}

// Sync flushes any buffered log entries to their respective destinations.
// It ensures that all pending log writes are completed before the program exits.
func Sync() {
	_ = Log.Sync()
}
