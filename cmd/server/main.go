package main

import (
	"fx-service/internal/application"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"fx-service/pkg/logger"
)

func main() {

	// Trim file paths to include only the relevant parts
	// Any file paths in error traces will be trimmed to the parts after "/fx-service"
	e.SetFilePathTrimPoint("/fx-service")

	// Set up a base logger
	appLogger := logger.NewLogger()

	// Create and boot the application
	app := application.NewApp(appLogger)

	if err := app.SetConfigs(); err != nil {
		appLogger.Error(err.Error(), nil)
		c.Out(err.Error())
		return
	}

	app.MonitorSignals().
		SetProviders().
		SetRoutes().
		Serve()
}
