package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/server"
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
	meHandler := NewMeHandler(s.ImageUrl)
	chatHandler := NewChat()
	uploadImageHandler := NewUploadImage(s.ImageUploader, s.ImageUrl)
	userHandler := NewUserHandler(s.ImageUrl)

	if s.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.Service(s.Service))
	router.Use(middleware.NewHandleErrorMiddleware(s.Conf.ServiceName))
	router.StaticFile("/", "template/index.html")

	api := router.Group("/v1/api")

	authRouter := api.Group("/auth")
	// 学生注册
	authRouter.POST("/register", authHandler.Register)
	// 登陆(老师、学生、管理员)
	authRouter.POST("/login", authHandler.Login)

	// uri: /v1/api/
	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware)
	{
		// 建立ws连接
		authorized.GET("/ws_conn", chatHandler.WsConn)
		// 上传图片
		authorized.POST("/upload_image", uploadImageHandler.UploadImage)
	}

	// logged uri: /v1/api/auth
	authRoute := authorized.Group("/auth")
	{
		authRoute.GET("/me", meHandler.Show)
		authRoute.GET("/logout", authHandler.Logout)
		// 老师列表
		authRoute.GET("/teacherList", userHandler.TeacherList)
	}

	// student uri: /v1/api/student
	student := authorized.Group("/student")
	student.Use(middleware.StudentMiddleware)
	{
	}

	// teacher uri: /v1/api/teacher
	teacher := authorized.Group("/teacher")
	teacher.Use(middleware.TeacherMiddleware)
	{
	}

	// admin uri: /v1/api/admin
	admin := authorized.Group("/admin")
	admin.Use(middleware.AdminMiddleware)
	{
	}

	return router
}
