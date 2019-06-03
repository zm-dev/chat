package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
)

type userHandler struct {
	imageUrl image_url.URL
}

func (u *userHandler) TeacherList(c *gin.Context) {
	users := make([]*model.User, 0, 4)
	users, err := service.TeacherList(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusOK, convert2UserListResp(c.Request.Context(), users, u.imageUrl))
}

func NewUserHandler(imageUrl image_url.URL) *userHandler {
	return &userHandler{imageUrl: imageUrl}
}
