package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghostmesh/vpnclient/internal/app"
	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/ghostmesh/vpnclient/internal/logger"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Initialize logger
	logger.Init(*debug)
	log := logger.GetLogger()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application instance
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Run the application
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}