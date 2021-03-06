package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/server"
)

func getInt32PageAndSize(c *gin.Context) (page, size int32) {
	var err error
	sizeI64, err := strconv.ParseInt(c.Query("size"), 10, 32)
	if err != nil {
		size = 10
	} else {
		size = int32(sizeI64)
	}
	if size > 50 {
		size = 50
	}
	pageI64, err := strconv.ParseInt(c.Query("page"), 10, 32)
	if err != nil {
		page = 0
	} else {
		page = int32(pageI64)
	}
	return page, size
}

func CreateHTTPHandler(s *server.Server) http.Handler {
	authHandler := NewAuthHandler()
	meHandler := NewMeHandler(s.ImageUrl)
	chatHandler := NewChatHandler()
	uploadImageHandler := NewUploadImage(s.ImageUploader, s.ImageUrl)
	userHandler := NewUserHandler(s.ImageUrl)
	recordHandler := NewRecordHandler(s.ImageUrl)

	if s.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.Service(s.Service))
	router.Use(middleware.NewHandleErrorMiddleware(s.Conf.ServiceName))

	router.Static("/assets", "./public")
	router.StaticFile("/", "public/index.html")

	api := router.Group("/api/v1")

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

		// 老师列表
		authorized.GET("teacher_list", userHandler.TeacherList)

		// 学生列表
		authorized.GET("student_list", userHandler.StudentList)

		// 聊天记录
		authorized.GET("/record", recordHandler.RecordListByUser)

		// 消息列表
		authorized.GET("/message_list", recordHandler.MessageList)

		// 批量设置消息为已读状态
		authorized.PUT("/record/batch_set_read", recordHandler.BatchSetRead)

		// 设置2人全部消息为已读
		authorized.PUT("/record/set_all_read", recordHandler.SetAllRead)

		// 显示指定用户信息
		authorized.GET("/user/:uid", userHandler.Show)

		// 不同类型的用户数量
		authorized.GET("user_count_by_type", userHandler.UserCountGroupByType)
	}

	// logged uri: /v1/api/auth
	authRoute := authorized.Group("/auth")
	{
		authRoute.GET("/me", meHandler.Show)
		authRoute.GET("/logout", authHandler.Logout)
		authRoute.PUT("/me", userHandler.UpdateMeInfo)
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
		admin.POST("/teacher", userHandler.CreateTeacher)
		admin.DELETE("/user/:uid", userHandler.DeleteUser)
		admin.GET("/record", recordHandler.AdminRecordList)
		admin.GET("/message_list", recordHandler.AdminMessageList)
		admin.PUT("/user/:uid", userHandler.UserUpdate)
	}

	return router
}
