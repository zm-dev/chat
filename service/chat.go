package service

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat/model"
	"sync"
	"time"
)

type ChatService struct {
	userWsConnMap sync.Map
}

func (c *ChatService) IsOnline(userId int64) bool {
	_, ok := c.userWsConnMap.Load(userId)
	return ok
}

func (c *ChatService) SendMsg(userId int64, msg model.IMsg) error {
	conn, ok := c.userWsConnMap.Load(userId);

	if !ok {
		return errors.New("用户不存在或不在线")
	}

	msg.SetSendAt(time.Now())
	return conn.(*websocket.Conn).WriteJSON(msg)
}

func (c *ChatService) OnLine(userId int64, conn *websocket.Conn) {
	c.userWsConnMap.Store(userId, conn)
}

func NewChatService() model.ChatService {
	return &ChatService{
		userWsConnMap: sync.Map{},
	}
}
