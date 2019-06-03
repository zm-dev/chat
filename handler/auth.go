package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authHandler struct{}

type Req struct {
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

func (authHandler) Login(c *gin.Context) {
	req := &Req{}
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
	c.JSON(http.StatusNoContent, nil)
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

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
