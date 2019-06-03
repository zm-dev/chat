package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat_v2/chat"
	"time"
)

type WsUser struct {
	Id   int64
	conn *websocket.Conn
}

func (u *WsUser) GetId() int64 {
	return u.Id
}

func (u *WsUser) SendMsg(msg chat.IMsg) error {
	return u.conn.WriteJSON(msg)
}

type Chat struct {
	piazza *chat.Piazza
}

type Input struct {
	Command    string `json:"command"`
	ChatRoomId string `json:"chat_room_id"`
	Msg        string `json:"msg"`
}

func (c *Chat) WsConn(ctx *gin.Context) {
	//userId, exists := ctx.Get("user_id")
	//if !exists {
	//	ctx.Error(errors.New("user_id 不存在"))
	//	return
	//}
	//userIdInt := userId.(int64)
	var userIdInt int64 = 5

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	defer conn.Close()
	if err != nil {
		ctx.Error(err)
		return
	}
	var input = new(Input)
	for {
		err := conn.ReadJSON(input)
		if err == nil {
			u := &WsUser{Id: userIdInt, conn: conn}
			switch input.Command {
			case "join":
				c.piazza.Join(input.ChatRoomId, u)
			case "leave":
				c.piazza.Leave(input.ChatRoomId, u)
			case "msg":
				c.piazza.Receive(input.ChatRoomId, &chat.Msg{User: u, Data: []byte(input.Msg), SendAt: time.Now()})
			}
		} else {
			fmt.Println(err)
		}
	}
}

func NewChat(piazza *chat.Piazza) Chat {
	return Chat{piazza: piazza}
}
