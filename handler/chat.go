package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
)

type Chat struct {
}

type Input struct {
	ToUserId int64  `json:"to_user_id"`
	Msg      string `json:"msg"`
}

func (c *Chat) WsConn(ctx *gin.Context) {
	userId, exists := ctx.Get(middleware.UserIdKey)
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
			// todo 错误处理
			fmt.Println(err)
			continue
		}
		msg := model.Msg{
			SendUserId: userIdInt,
			Data:       []byte(input.Msg),
		}
		// websocket 发送消息
		err = service.SendMsg(ctx.Request.Context(), input.ToUserId, &msg)
		if err != nil {
			// todo 错误处理
			fmt.Println(err)
		}

		err = service.CreateRecord(ctx.Request.Context(), &model.Record{
			FromId: userIdInt,
			ToId:   input.ToUserId,
			Msg:    input.Msg,
		})

		if err != nil {
			// todo 错误处理
			fmt.Println(err)
		}
	}
}

func NewChat() Chat {
	return Chat{}
}
