package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat_v2/errors"
	"github.com/zm-dev/chat_v2/model"
	"github.com/zm-dev/chat_v2/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHandler struct{}

func (authHandler) Login(c *gin.Context) {
	req := &struct {
		Account  string `form:"account" json:"account"`
		Password string `form:"password" json:"password"`
	}{}
	if err := c.ShouldBind(req); err != nil {
		c.Error(errors.BindError(err))
		return
	}
	ticket, err := service.UserLogin(c.Request.Context(), strings.TrimSpace(req.Account), strings.TrimSpace(req.Password))
	if err != nil {
		c.Error(err)
		return
	}
	setAuthCookie(c, ticket.Id, ticket.UserId, int(ticket.ExpiredAt.Sub(time.Now()).Seconds()))
	c.JSON(http.StatusNoContent, nil)
}

func (authHandler) Logout(c *gin.Context) {
	ticketId, err := c.Cookie("ticket_id")
	if err != nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	removeAuthCookie(c)
	service.TicketDestroy(c.Request.Context(), ticketId)
	c.JSON(http.StatusNoContent, nil)
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
func (authHandler) Register(c *gin.Context) {
	type Req struct {
		Account  string `form:"account" json:"account"`
		Password string `form:"password" json:"password"`
	}
	req := &Req{}
	if err := c.ShouldBind(req); err != nil {
		_ = c.Error(err)
		return
	}
	_, err := service.UserRegister(c.Request.Context(), strings.TrimSpace(req.Account), model.CertificateType(0), req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(201)
}

func setAuthCookie(c *gin.Context, ticketId string, userId int64, maxAge int) {
	c.SetCookie("ticket_id", ticketId, maxAge, "", "", false, true)
	c.SetCookie("user_id", strconv.FormatInt(userId, 10), maxAge, "", "", false, false)
}

func removeAuthCookie(c *gin.Context) {
	c.SetCookie("ticket_id", "", -1, "", "", false, true)
	c.SetCookie("user_id", "", -1, "", "", false, false)
}
