package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHandler struct{}

type RegisterReq struct {
	Nickname string `form:"nickname" json:"nickname"`
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

type LoginReq struct {
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

func (authHandler) Login(c *gin.Context) {
	req := &LoginReq{}

	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}

	ticket, err := service.UserLogin(c.Request.Context(), strings.TrimSpace(req.Account), strings.TrimSpace(req.Password))
	if err != nil {
		_ = c.Error(err)
		return
	}
	setAuthCookie(c, ticket.Id, ticket.UserId, int(ticket.ExpiredAt.Sub(time.Now()).Seconds()))
	u, err := service.UserLoad(c.Request.Context(), ticket.UserId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	certificate, err := service.CertificateLoadByUserId(c.Request.Context(), ticket.UserId)
	if err != nil {
		c.Abort()
		return
	}
	userType := enum.CertificateStudent
	if certificate.Type == enum.CertificateTeacher {
		userType = enum.CertificateTeacher
	} else if certificate.Type == enum.CertificateAdmin {
		userType = enum.CertificateAdmin
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"userType": userType,
		"userId":   u.Id,
	})
}

func (authHandler) Logout(c *gin.Context) {
	ticketId, err := c.Cookie("ticket_id")
	if err != nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	removeAuthCookie(c)
	_ = service.TicketDestroy(c.Request.Context(), ticketId)
	c.JSON(http.StatusNoContent, nil)
}

func (authHandler) Register(c *gin.Context) {
	req := &RegisterReq{}
	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(err)
		return
	}
	_, err := service.UserRegister(c.Request.Context(), strings.TrimSpace(req.Account), strings.TrimSpace(req.Password), req.Nickname, enum.CertificateStudent)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 登陆
	ticket, err := service.UserLogin(c.Request.Context(), strings.TrimSpace(req.Account), strings.TrimSpace(req.Password))
	if err != nil {
		_ = c.Error(err)
		return
	}
	setAuthCookie(c, ticket.Id, ticket.UserId, int(ticket.ExpiredAt.Sub(time.Now()).Seconds()))
	c.Status(http.StatusCreated)
}

func setAuthCookie(c *gin.Context, ticketId string, userId int64, maxAge int) {
	////TODO 临时添加一条
	//c.SetCookie("ticket_id", ticketId, maxAge, "", "dev.hn-zm.com", false, true)
	//c.SetCookie("user_id", strconv.FormatInt(userId, 10), maxAge, "dev.hn-zm.com", "", false, false)

	c.SetCookie("ticket_id", ticketId, maxAge, "", "", false, true)
	c.SetCookie("user_id", strconv.FormatInt(userId, 10), maxAge, "", "", false, false)
}

func removeAuthCookie(c *gin.Context) {
	//c.SetCookie("ticket_id", "", -1, "", "dev.hn-zm.com", false, true)
	//c.SetCookie("user_id", "", -1, "", "dev.hn-zm.com", false, false)

	c.SetCookie("ticket_id", "", -1, "", "", false, true)
	c.SetCookie("user_id", "", -1, "", "", false, false)
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
