package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
	"sort"
	"time"
)

type recordHandler struct {
	imageUrl image_url.URL
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
	NickName            string    `json:"nick_name"`
	LastMessage         string    `json:"last_message"`
	IsMeSend            bool      `json:"is_me_send"`
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
	userIdInt := middleware.UserId(ctx)

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
	authId := middleware.UserId(c)
	size := 20
	records, err := service.LastRecordList(c.Request.Context(), authId, size)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result := make([]*MessageRecordListResponse, 0, len(records))
	for _, item := range records {
		itemMsg := &MessageRecordListResponse{}

		// 对方用户的信息（我发消息给你 OR 你发消息给我，这里都显示你的基本信息）
		user := new(model.User)
		if item.FromId == authId {
			itemMsg.IsMeSend = true
			user, err = service.UserLoad(c.Request.Context(), item.ToId)
		} else {
			user, err = service.UserLoad(c.Request.Context(), item.FromId)
		}
		if err != nil || user.Id == 0 {
			_ = c.Error(errors.BadRequest("用户不存在"))
			return
		}
		itemMsg.UserId = user.Id
		itemMsg.NickName = user.NickName
		itemMsg.AvatarUrl = r.imageUrl.Generate(user.AvatarHash)

		// 消息
		itemMsg.LastMessage = item.Msg
		itemMsg.LastMessageSendTime = item.CreatedAt
		itemMsg.NotReadMsgCount = service.GetNotReadRecordCount(c.Request.Context(), itemMsg.UserId, authId)

		result = append(result, itemMsg)
	}
	// 输出结果排序
	sort.Slice(result, func(i, j int) bool {
		return /*result[i].NotReadMsgCount <= 0 &&*/ result[i].LastMessageSendTime.After(result[j].LastMessageSendTime)
	})
	c.JSON(http.StatusOK, result)
}

func NewRecordHandler(imageUrl image_url.URL) *recordHandler {
	return &recordHandler{imageUrl: imageUrl}
}
