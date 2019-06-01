package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat_v2/handler/middleware"
	"github.com/zm-dev/chat_v2/server"
	"net/http"
	"strconv"
)

func getInt32LimitAndOffset(c *gin.Context) (limit, offset int32) {
	var err error
	limitI64, err := strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 10
	} else {
		limit = int32(limitI64)
	}
	if limit > 50 {
		limit = 50
	}

	offsetI64, err := strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	} else {
		offset = int32(offsetI64)
	}
	return limit, offset
}

func getInt64LimitAndOffset(c *gin.Context) (limit, offset int64) {
	var err error
	limit, err = strconv.ParseInt(c.Query("limit"), 10, 32)
	if err != nil {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	offset, err = strconv.ParseInt(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}
	return limit, offset
}

func CreateHTTPHandler(s *server.Server) http.Handler {
	authHandler := NewAuthHandler()
	meHandler := NewMeHandler()
	if s.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.Gorm(s.DB))
	router.Use(middleware.Service(s.Service))
	router.Use(middleware.NewHandleErrorMiddleware(s.Conf.ServiceName))
	api := router.Group("/api")
	authRouter := api.Group("/auth")
	{
		authRouter.POST("/login", authHandler.Login)
		authRouter.GET("/logout", authHandler.Logout)
		authRouter.POST("/register", authHandler.Register)
	}

	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware)
	{
		authorized.GET("/me", meHandler.Show)
	}
	adminRouter := api.Group("/")
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	{
	}
	return router
}
