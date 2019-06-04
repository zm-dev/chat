package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
	"time"
)

type recordHandler struct {
}

type RecordListRequest struct {
	UserIdB int64 `form:"user_id_b" json:"user_id_b"`
}

type BatchSetReadRequest struct {
	RecordIds []int64 `form:"record_ids" json:"record_ids"` // 账号
}

// 消息列表
type MessageRecordListResponse struct {
	UserId              int64     `json:"user_id"`
	AvatarUrl           string    `json:"avatar_url"`
	NikeName            string    `json:"nike_name"`
	LastMessage         string    `json:"last_message"`
	NotReadMsgCount     int32     `json:"not_read_msg_count"`
	LastMessageSendTime time.Time `json:"last_message_send_time"`
}

func (r *recordHandler) RecordListByUser(c *gin.Context) {
	req := &RecordListRequest{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	page, size := getInt32PageAndSize(c)
	userId := middleware.UserId(c)
	result := &model.Page{
		Current: page,
		Size:    size,
	}
	err := service.PageRecord(c.Request.Context(), result, userId, req.UserIdB, false)
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusOK, result)
}

func (r *recordHandler) BatchSetRead(ctx *gin.Context) {
	userId, exists := ctx.Get(middleware.UserIdKey)
	if !exists {
		_ = ctx.Error(errors.ErrAccountNotFound())
		return
	}
	userIdInt := userId.(int64)

	req := &BatchSetReadRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}

	err := service.BatchSetRead(ctx.Request.Context(), req.RecordIds, userIdInt);
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// 与当前用户聊过天的用户列表以及最近的一条消息
// 最近聊过天的用户列表（如果未读消息超过20条则全部显示，否则显示20条有未读消息和没有未读消息的用户列表）
func (r *recordHandler) MessageList(c *gin.Context) {

}

func NewRecordHandler() *recordHandler {
	return &recordHandler{}
}
