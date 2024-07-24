package ginHandlers

import (
	"fx-service/internal/middleware"
	"fx-service/pkg/config"
	"fx-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Logger *logger.Logger
	Config *config.Config
	Engine *gin.Engine
}

func (r *GinRouter) RegisterMiddleware() {
	r.Engine.Use(middleware.GinLogger(r.Logger))
	r.Engine.Use(middleware.GinRateLimiter(r.Config.RateLimiter))

	// Set default content type to JSON
	r.Engine.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	})
}

func (r *GinRouter) RegisterRoutes() {
	r.Engine.GET("/rate/:from/:to", GetRate(r.Config))
	r.Engine.GET("/rates", GetRates(r.Config))
	r.Engine.GET("/status", GetStatus(r.Config))
}

func (r *GinRouter) Serve(addr string) error {
	return r.Engine.Run(addr)
}
