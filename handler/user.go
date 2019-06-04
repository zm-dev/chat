package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat/enum"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"net/http"
)

type userHandler struct {
	imageUrl image_url.URL
}

type UserCreateRequest struct {
	Account    string `form:"account" json:"account"`         // 账号
	Password   string `form:"password" json:"password"`       // 密码
	AvatarHash string `form:"avatar_hash" json:"avatar_hash"` // 头像
	NikeName   string `form:"nike_name" json:"nike_name"`     // 昵称
	Profile    string `form:"profile" json:"profile"`         // 简介
	Company    string `form:"company" json:"company"`         // 工作单位
	Gender     uint8  `form:"gender" json:"gender"`           // 性别
	GroupId    uint8  `form:"groupId" json:"groupId"`         // 组
}

func (u *userHandler) TeacherList(c *gin.Context) {
	users := make([]*model.User, 0, 4)
	users, err := service.TeacherList(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusOK, convert2UserListResp(c.Request.Context(), users, u.imageUrl))
}

func (u *userHandler) CreateTeacher(c *gin.Context) {
	req := &UserCreateRequest{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	// 老师密码不传默认为账号
	if req.Password == "" {
		req.Password = req.Account
	}
	userId, err := service.UserRegister(c.Request.Context(), req.Account, enum.CertificateTeacher, req.Password)
	if err != nil {
		_ = c.Error(err)
	}
	err = service.UserUpdate(c.Request.Context(), &model.User{
		Id:         userId,
		AvatarHash: req.AvatarHash,
		NikeName:   req.NikeName,
		Profile:    req.Profile,
		Gender:     enum.Gender(req.Gender),
		GroupId:    enum.Group(req.GroupId),
		Company:    req.Company,
	})
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusNoContent, nil)
}

func NewUserHandler(imageUrl image_url.URL) *userHandler {
	return &userHandler{imageUrl: imageUrl}
}
