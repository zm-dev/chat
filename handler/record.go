package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/service"
	"net/http"
)

type Record struct {
}

type BatchSetReadRequest struct {
	RecordIds []int64 `form:"record_ids" json:"record_ids"` // 账号
}

func (c *Record) BatchSetRead(ctx *gin.Context) {
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

func NewRecordHandler() Record {
	return Record{}
}
