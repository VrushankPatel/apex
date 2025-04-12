package util

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// InitLogger sets up the logger with the desired configuration
func InitLogger() {
	// Ensure the data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory: ", err)
	}

	// Create log file
	logFile, err := os.OpenFile(
		filepath.Join("data", "arbitrage.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	// Configure logrus to write to both file and stdout
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// Set log level to info by default
	log.SetLevel(log.InfoLevel)

	// Set up formatter
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.Info("Logger initialized")
}
