package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type IMsg interface {
	GetUserId() int64
	GetData() []byte
	GetSendAt() time.Time
	SetSendAt(time.Time)
	MarshalJSON() ([]byte, error)
}

type Msg struct {
	SendUserId int64
	Data       []byte
	SendAt     time.Time
}

func (m *Msg) GetUserId() int64 {
	if m == nil {
		return -1
	}
	return m.GetUserId()
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

func (m *Msg) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		UserId int64     `json:"user_id"`
		Data   string    `json:"data"`
		SendAt time.Time `json:"send_at"`
	}{
		UserId: m.GetUserId(),
		Data:   string(m.Data),
		SendAt: m.SendAt,
	})
}

type ChatService interface {
	// 用户是否在线
	IsOnline(userId int64) bool
	SendMsg(userId int64, msg IMsg) error
	// 用户上线
	OnLine(userId int64, conn *websocket.Conn)
}
