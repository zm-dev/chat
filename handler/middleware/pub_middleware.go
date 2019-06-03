package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/queue"
)

func Pub(pub queue.PubQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(queue.NewContext(c.Request.Context(), pub))
		c.Next()
	}
}
