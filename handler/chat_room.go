package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zm-dev/chat_v2/chat"
	"github.com/zm-dev/go-image_uploader/image_url"
	"strings"
	"time"
)

type chatRoomResp struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CoverHash string `json:"cover_hash"`
	CoverUrl  string `json:"cover_url"`
}

func convert2chatRoomResp(chatRoom *model.ChatRoom, imageUrl image_url.URL) *chatRoomResp {
	return &chatRoomResp{
		Id:        chatRoom.Id,
		Name:      chatRoom.Name,
		CoverHash: chatRoom.CoverHash,
		CoverUrl:  imageUrl.Generate(chatRoom.CoverHash, image_url.Height(200), image_url.Width(200)),
	}
}

func covert2chatRoomRespList(chatRooms []*model.ChatRoom, imageUrl image_url.URL) []*chatRoomResp {
	chatRoomResps := make([]*chatRoomResp, 0, len(chatRooms))
	for _, v := range chatRooms {
		chatRoomResps = append(chatRoomResps, convert2chatRoomResp(v, imageUrl))
	}
	return chatRoomResps
}

type ChatRoom struct {
	piazza   *chat.Piazza
	imageUrl image_url.URL
}

func (cr *ChatRoom) List(c *gin.Context) {
	chatRooms := middleware.GetService(c).ChatRoomList()
	c.JSON(200, covert2chatRoomRespList(chatRooms, cr.imageUrl))
}
func (cr *ChatRoom) Create(c *gin.Context) {
	type Req struct {
		Name      string `json:"name" form:"name"`
		CoverHash string `json:"cover_hash" form:"cover_hash"`
	}

	req := &Req{}

	if err := c.ShouldBind(req); err != nil {
		c.Error(errors.New("参数有误"))
		return
	}

	chatRoomModel := &model.ChatRoom{
		Name:      req.Name,
		CoverHash: req.CoverHash,
	}
	err := middleware.GetService(c).ChatRoomCreate(chatRoomModel)

	chatRoom := cr.piazza.CreateChatRoom(chatRoomModel.Id)
	go chatRoom.Run()

	chatRoom.OnJoin(func(user chat.User) {
		u, err := middleware.GetService(c).UserLoad(5)
		if err == nil {
			chatRoom.Broadcast(&chat.Msg{
				User:   &chat.SysUser{},
				Data:   []byte(fmt.Sprintf("%s 进入房间", u.NikeName)),
				SendAt: time.Now(),
			})
		}
	})

	chatRoom.OnLeave(func(user chat.User) {
		u, err := middleware.GetService(c).UserLoad(user.GetId())
		if err == nil {
			chatRoom.Broadcast(&chat.Msg{
				User:   &chat.SysUser{},
				Data:   []byte(fmt.Sprintf("%s 已经离开", u.NikeName)),
				SendAt: time.Now(),
			})
		}
	})

	chatRoom.ProcessMsg(func(msg chat.IMsg) {
		uid := msg.GetUser().GetId()
		user, err := middleware.GetService(c).UserLoad(uid)
		if err == nil {
			// userResp := convert2UserResp(user, cr.imageUrl)
			b, err := json.Marshal(user)
			if err == nil {
				msg.SetMeta(map[string]string{
					"user": string(b),
				})
			}
		}
	})

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(204, nil)
}

func (cr *ChatRoom) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.Error(errors.New("参数有误"))
		return
	}
	err := middleware.GetService(c).ChatRoomDelete(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(204, nil)
}

func NewChatRoom(piazza *chat.Piazza, imageUrl image_url.URL) ChatRoom {
	return ChatRoom{piazza, imageUrl}
}
