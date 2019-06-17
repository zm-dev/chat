package service

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/zm-dev/chat/model"
	"log"
	"strconv"
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

	c.Broadcast(&model.Msg{
		SendUserId: -1,
		SendAt:     time.Now(),
		Meta: map[string]string{
			"msg_type": "user_online",
			"user_id":  strconv.FormatInt(userId, 10),
		},
	}, []int64{userId})
	log.Printf("[OnLine] %d 用户上线啦\n", userId)
	c.userWsConnMap.Store(userId, conn)
}

func (c *ChatService) OffLine(userId int64) {

	c.Broadcast(&model.Msg{
		SendUserId: -1,
		SendAt:     time.Now(),
		Meta: map[string]string{
			"msg_type": "user_offline",
			"user_id":  strconv.FormatInt(userId, 10),
		},
	}, []int64{userId})
	log.Printf("[OffLine] %d 用户退出啦\n", userId)
	c.userWsConnMap.Delete(userId)
}

func (c *ChatService) Broadcast(msg model.IMsg, excludeUserId []int64) {
	c.userWsConnMap.Range(func(key, value interface{}) bool {
		for v := range excludeUserId {
			if key == v {
				return true
			}
		}
		if value != nil {
			_ = value.(*websocket.Conn).WriteJSON(msg)
		}
		return true
	})
}

func IsOnline(c context.Context, userId int64) bool {
	return FromContext(c).IsOnline(userId)
}

func SendMsg(c context.Context, userId int64, msg model.IMsg) error {
	return FromContext(c).SendMsg(userId, msg)
}

func OnLine(c context.Context, userId int64, conn *websocket.Conn) {
	FromContext(c).OnLine(userId, conn)
}

func OffLine(c context.Context, userId int64) {
	FromContext(c).OffLine(userId)
}

func Broadcast(c context.Context, msg model.IMsg, excludeUserId []int64) {
	FromContext(c).Broadcast(msg, excludeUserId)
}

func NewChatService() model.ChatService {
	return &ChatService{
		userWsConnMap: sync.Map{},
	}
}
