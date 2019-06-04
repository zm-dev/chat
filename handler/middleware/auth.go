package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
)

var (
	IsLoginKey    = "is_login"
	UserIdKey     = "user_id"
	LoggedUserKey = "logged_user"
)

func AuthMiddleware(c *gin.Context) {
	isLogin := check(c)
	if !isLogin {
		_ = c.Error(errors.Unauthorized())
		c.Abort()
		return
	}
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	isLogin := check(c)
	if !isLogin {
		_ = c.Error(errors.Unauthorized())
		c.Abort()
		return
	}
	uId := UserId(c)
	certificate, err := service.CertificateLoadByUserId(c.Request.Context(), uId)
	if err != nil {
		c.Abort()
		return
	}
	if certificate.Type != enum.CertificateAdmin {
		_ = c.Error(errors.Forbidden("没有权限"))
		c.Abort()
		return
	}
	c.Next()
}

func TeacherMiddleware(c *gin.Context) {
	isLogin := check(c)
	if !isLogin {
		_ = c.Error(errors.Unauthorized())
		c.Abort()
		return
	}
	uId := UserId(c)
	certificate, err := service.CertificateLoadByUserId(c.Request.Context(), uId)
	if err != nil {
		c.Abort()
		return
	}
	if certificate.Type != enum.CertificateTeacher {
		_ = c.Error(errors.Forbidden("没有权限"))
		c.Abort()
		return
	}
	c.Next()
}

func StudentMiddleware(c *gin.Context) {
	isLogin := check(c)
	if !isLogin {
		_ = c.Error(errors.Unauthorized())
		c.Abort()
		return
	}
	uId := UserId(c)
	certificate, err := service.CertificateLoadByUserId(c.Request.Context(), uId)
	if err != nil {
		c.Abort()
		return
	}
	if certificate.Type != enum.CertificateStudent {
		_ = c.Error(errors.Forbidden("没有权限"))
		c.Abort()
		return
	}
	c.Next()
}

func check(c *gin.Context) bool {
	var (
		isLogin bool
	)
	if ticketId, err := c.Cookie("ticket_id"); err == nil {
		isValid, userId, err := service.TicketIsValid(c.Request.Context(), ticketId)
		if err == nil {
			isLogin = isValid
			setIsLogin(c, isLogin)
			setUserId(c, userId)
		}
	} else {
		// cookie不存在
		isLogin = false
	}
	return isLogin
}

func setIsLogin(c *gin.Context, isLogin bool) {
	c.Set(IsLoginKey, isLogin)
}

func setUserId(c *gin.Context, userId int64) {
	c.Set(UserIdKey, userId)
}

func CheckLogin(c *gin.Context) bool {
	isLogin, ok := c.Get(IsLoginKey)
	if !ok {
		return check(c)
	}
	return isLogin.(bool)

}

func UserId(c *gin.Context) int64 {
	userId, ok := c.Get(UserIdKey)
	if !ok {
		check(c)
		return c.GetInt64(UserIdKey)
	}
	return userId.(int64)
}

func LoggedUser(c *gin.Context) *model.User {
	user, ok := c.Get(LoggedUserKey)
	if !ok {
		userId := UserId(c)
		if userId == 0 {
			return nil
		}
		userModel, err := service.UserLoad(c.Request.Context(), userId)
		if err != nil {
			return nil
		}
		c.Set("LoggedUserKey", userModel)
		return userModel
	}
	return user.(*model.User)
}
