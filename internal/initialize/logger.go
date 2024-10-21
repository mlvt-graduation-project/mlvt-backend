package initialize

import (
	"fmt"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/infra/zap-logging/zap"
	"os"
	"path/filepath"
)

// InitLogger sets up the logging system based on the loaded configuration.
func InitLogger() error {
	logLevel := env.EnvConfig.LogLevel
	logPath := env.EnvConfig.LogPath

	// Ensure the log directory exists
	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// Initialize logging
	log.SetLogger(zap.NewLogger(
		log.ParseLevel(logLevel),
		zap.WithName("mlvt"),
		zap.WithPath(logPath),
		zap.WithCallerFullPath(),
	))

	return nil
}
