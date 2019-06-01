package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat_v2/store/db_store"
)

func Gorm(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(db_store.NewDBContext(c.Request.Context(), db))
		c.Next()
	}
}
