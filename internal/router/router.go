package router

import (
	"fx-service/internal/router/fiber"
	"fx-service/internal/router/gin"
	"fx-service/pkg/config"
	"fx-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	RegisterMiddleware()
	RegisterRoutes()
	Serve(addr string) error
}

// NewFiberRouter creates a new instance of the FiberRouter
func NewFiberRouter(logger *logger.Logger, config *config.Config) Router {
	// Create a new Fiber app with a custom error handler
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	return &fiberHandlers.FiberRouter{
		Logger: logger,
		Config: config,
		App:    fiberApp,
	}
}

// NewGinRouter creates a new instance of the GinRouter
func NewGinRouter(logger *logger.Logger, config *config.Config) Router {
	//gin.SetMode(gin.ReleaseMode) // TODO - Set this in the config
	engine := gin.New()

	return &ginHandlers.GinRouter{
		Logger: logger,
		Config: config,
		Engine: engine,
	}
}
