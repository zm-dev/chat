package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat_v2/handler/middleware"
	"github.com/zm-dev/chat_v2/model"
	"github.com/zm-dev/chat_v2/service"
	"github.com/zm-dev/chat_v2/util"
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
	c.JSON(http.StatusOK, convert2UserResp(user, m.imageUrl))
}

func convert2UserListResp(users []*model.User, imageUrl image_url.URL) []map[string]interface{} {
	userList := make([]map[string]interface{}, 0, len(users))
	for _, v := range users {
		userList = append(userList, convert2UserResp(v, imageUrl))
	}
	return userList
}

func convert2UserResp(user *model.User, imageUrl image_url.URL) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.Id,
		"name":       user.NikeName,
		"email":      user.Email,
		"avatarUrl":  imageUrl.Generate(user.AvatarHash),
		"profile":    user.Profile,
		"gender":     util.ConvertUserGender(user.Gender),
		"status":     util.ConvertUserStatus(user.Status),
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}

func NewMeHandler(imageUrl image_url.URL) *meHandler {
	return &meHandler{imageUrl: imageUrl}
}
