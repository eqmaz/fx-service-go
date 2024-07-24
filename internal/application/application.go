package application

import (
	"fx-service/internal/router"
	"fx-service/internal/service/providers"
	"fx-service/internal/service/ratecache"
	"fx-service/pkg/config"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"fx-service/pkg/logger"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type App struct {
	Config *config.Config
	Logger *logger.Logger
	Router router.Router
}

// NewApp initializes the App with the necessary dependencies
func NewApp(logger *logger.Logger) *App {
	return &App{
		Logger: logger,
	}
}

// SetConfigs loads the application configuration from the defaults, config.json file and environment variables
func (app *App) SetConfigs() error {
	// Set up the error catalogue
	e.SetCatalogue(errorMap)

	// Try to locate the required config file.
	// It will first look for arguments passed in, then it will look in the current working dir and the executable dir.
	// If we don't have one, the default values will be used.
	configFilePath, err := getConfigFilePath("config.json")

	// Load the config
	appConfig, err := config.NewConfig(configFilePath)
	if err != nil {
		return err
	}

	// Convert the supported currencies to uppercase for consistency
	appConfig.CurrenciesToUppercase()

	ratecache.GetInstance().SetExpiry(appConfig.CacheExpirySec)

	app.Config = &appConfig

	return nil
}

// SetProviders initializes API providers by checking API keys and loading supported currencies for enabled providers
func (app *App) SetProviders() *App {
	// Ensure the Config is set
	if app.Config == nil {
		c.Out("Config not set. Cannot initialize providers.")
		return app
	}

	// Initialize the providers - sets up API keys, etc.
	providers.InitProviders(&app.Config.Providers, app.Config.APITimeout)

	return app
}

// SetRoutes initializes the Router, route handlers and middleware
func (app *App) SetRoutes() *App {
	routerChoice := strings.ToLower(app.Config.Router)
	switch routerChoice {
	case "gin":
		app.Router = router.NewGinRouter(app.Logger, app.Config)
		c.Info("Using Gin router")
	default:
		app.Router = router.NewFiberRouter(app.Logger, app.Config)
		c.Info("Using Fiber router")
	}

	app.Router.RegisterMiddleware()
	app.Router.RegisterRoutes()
	return app
}

// Serve starts listening to requests on the specified port
func (app *App) Serve() {
	portStr := strconv.FormatUint(app.Config.Port, 10)
	c.Out("Application started on port " + portStr)
	c.Outf(" -> Running in '%s' mode", app.Config.Mode.String())

	address := ":" + portStr // e.g. :8080
	err := app.Router.Serve(address)
	if err != nil {
		c.Warnf("Server could not listen on the given port")
		e.FromError(err).Print(0, 0)
		return
	}
}

// MonitorSignals listens for OS signals and exits the program gracefully
func (app *App) MonitorSignals() *App {
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Notify the channel on receiving SIGINT (Ctrl+C) or SIGTERM
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to notify the program to exit
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		c.Out("Received: ", sig)
		done <- true
	}()

	go func() {
		<-done
		c.Out("Stopping server...")
		//app.Router.Stop()
		os.Exit(0)
	}()

	//c.Out("Press Ctrl+C to exit")

	return app
}

func GetErrorMap() map[string]string {
	return errorMap
}
