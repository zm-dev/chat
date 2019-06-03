package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
)

type meHandler struct {
	imageUrl image_url.URL
}

func (m *meHandler) Show(c *gin.Context) {
	uid := middleware.UserId(c)
	user, err := service.UserLoad(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, convert2UserResp(c.Request.Context(), user, m.imageUrl))
}

func convert2UserListResp(c context.Context, users []*model.User, imageUrl image_url.URL) []map[string]interface{} {
	userList := make([]map[string]interface{}, 0, len(users))
	for _, v := range users {
		userList = append(userList, convert2UserResp(c, v, imageUrl))
	}
	return userList
}

func convert2UserResp(c context.Context, user *model.User, imageUrl image_url.URL) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.Id,
		"name":       user.NikeName,
		"email":      user.Email,
		"avatarUrl":  imageUrl.Generate(user.AvatarHash),
		"profile":    user.Profile,
		"gender":     enum.ParseGender(user.Gender),
		"group":      enum.ParseGroup(user.GroupId),
		"status":     service.IsOnline(c, user.Id),
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}

func NewMeHandler(imageUrl image_url.URL) *meHandler {
	return &meHandler{imageUrl: imageUrl}
}
