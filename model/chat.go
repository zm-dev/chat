package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type IMsg interface {
	GetSendUserId() int64
	GetData() []byte
	GetSendAt() time.Time
	SetSendAt(time.Time)
	GetMeta() map[string]string
	SetMeta(map[string]string)
	MarshalJSON() ([]byte, error)
}

type Msg struct {
	SendUserId int64
	Meta       map[string]string
	Data       []byte
	SendAt     time.Time
}

func (m *Msg) GetSendUserId() int64 {
	if m == nil {
		return -1
	}
	return m.SendUserId
}

func (m *Msg) GetData() []byte {
	if m == nil {
		return nil
	}
	return m.Data
}

func (m *Msg) GetSendAt() time.Time {
	if m == nil {
		return time.Time{}
	}
	return m.SendAt
}

func (m *Msg) SetSendAt(t time.Time) {
	m.SendAt = t
}

func (m *Msg) GetMeta() map[string]string {
	return m.Meta
}

func (m *Msg) SetMeta(meta map[string]string) {
	m.Meta = meta
}

func (m *Msg) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		SendUserId int64             `json:"send_user_id"`
		Data       string            `json:"data"`
		SendAt     time.Time         `json:"send_at"`
		Meta       map[string]string `json:"meta"`
	}{
		SendUserId: m.GetSendUserId(),
		Data:       string(m.Data),
		SendAt:     m.SendAt,
		Meta:       m.GetMeta(),
	})
}

type ChatService interface {
	// 用户是否在线
	IsOnline(userId int64) bool
	SendMsg(userId int64, msg IMsg) error
	// 用户上线
	OnLine(userId int64, conn *websocket.Conn)
	// 用户下线
	OffLine(userId int64)
	// 广播消息
	Broadcast(msg IMsg, excludeUserId []int64)
}
