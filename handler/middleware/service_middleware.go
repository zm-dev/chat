package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat_v2/service"
)

func Service(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(service.NewContext(c.Request.Context(), svc))
		c.Next()
	}
}
