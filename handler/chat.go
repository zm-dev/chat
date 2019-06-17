package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat/errors"
	"github.com/zm-dev/chat/handler/middleware"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"strconv"
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
		_ = ctx.Error(errors.ErrAccountNotFound())
		return
	}
	userIdInt := userId.(int64)

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	// 用户上线
	service.OnLine(ctx.Request.Context(), userIdInt, conn)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var input = new(Input)

	for {
		err := conn.ReadJSON(input)
		if err != nil {
			// 调用用户下线
			service.OffLine(ctx.Request.Context(), userIdInt)
			break
		}


		recordId, err := service.CreateRecord(ctx.Request.Context(), &model.Record{
			FromId: userIdInt,
			ToId:   input.ToUserId,
			Msg:    input.Msg,
		})

		if err != nil {
			// 创建消息失败
			continue
		}

		msg := model.Msg{
			SendUserId: userIdInt,
			Data:       []byte(input.Msg),
			Meta: map[string]string{
				"msg_type": "record",
				"record_id":  strconv.FormatInt(recordId, 10),
			},
		}

		// websocket 发送消息
		err = service.SendMsg(ctx.Request.Context(), input.ToUserId, &msg)
		if err != nil {
			// todo 错误处理
			fmt.Println(err)
		}

	}

	conn.Close()
}

func NewChatHandler() Chat {
	return Chat{}
}
