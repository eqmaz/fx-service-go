package ginHandlers

import (
	"fx-service/internal/reply"
	"github.com/gin-gonic/gin"
	"net/http"
)

// replyResult sends a successful response to the client
func replyResult(c *gin.Context, data interface{}) {
	payload := reply.Result(data)
	c.JSON(http.StatusOK, payload)
}

// replyError sends an error response to the client, with the given status code
func replyError(c *gin.Context, status int, data interface{}) {
	payload := reply.Error(data)
	c.JSON(status, payload)
}
