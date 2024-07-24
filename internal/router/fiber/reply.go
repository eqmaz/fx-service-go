package fiberHandlers

import (
	"fx-service/internal/reply"
	"fx-service/internal/service/stats"
	"fx-service/pkg/e"
	"github.com/gofiber/fiber/v2"
)

// replyResult sends a successful response to the client
func replyResult(c *fiber.Ctx, data interface{}) error {
	c.Status(fiber.StatusOK)
	c.Set("content-type", "application/json; charset=utf-8")
	defer stats.GetInstance().IncRequestCount()
	payload := reply.Result(data)
	return c.JSON(payload)
}

// replyError sends an error response to the client
func replyError(c *fiber.Ctx, status int, data interface{}) error {
	var payload interface{}
	c.Status(status)
	//c.Set("content-type", "application/json; charset=utf-8")

	// if date is of type string
	switch data.(type) {
	case e.Exception:
		exception := data.(e.Exception)
		msg := exception.GetMessage()
		code := exception.GetCode()
		payload = reply.Error(map[string]string{"code": code, "message": msg})
	default:
		payload = reply.Error(data)
	}

	defer stats.GetInstance().IncRequestCount()
	return c.JSON(payload)
}
