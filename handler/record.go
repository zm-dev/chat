package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
)

type recordHandler struct {
}

type RecordListRequest struct {
	UserIdB int64 `form:"user_id_b" json:"user_id_b"`
}

func (r *recordHandler) RecordList(c *gin.Context) {
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

type BatchSetReadRequest struct {
	RecordIds []int64 `form:"record_ids" json:"record_ids"` // 账号
}

func (c *recordHandler) BatchSetRead(ctx *gin.Context) {
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

func NewRecordHandler() *recordHandler {
	return &recordHandler{}
}
