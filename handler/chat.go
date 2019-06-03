package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat/model"
)

type Chat struct {
	chatService model.ChatService
}

type Input struct {
	SendUserId string `json:"send_user_id"`
	Msg        string `json:"msg"`
}

func (c *Chat) WsConn(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		_ = ctx.Error(errors.New("user_id 不存在"))
		return
	}
	userIdInt := userId.(int64)

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	defer conn.Close()
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	var input = new(Input)
	for {
		err := conn.ReadJSON(input)
		if err != nil {
			// todo 返回错误
			fmt.Println(err)
		}
		msg := model.Msg{
			SendUserId: userIdInt,
			Data:       []byte(input.Msg),
		}
		err = c.chatService.SendMsg(userIdInt, &msg)
		if err != nil {
			// todo 返回错误
			fmt.Println(err)
		}
	}
}

func NewChat(chatService model.ChatService) Chat {
	return Chat{chatService}
}
