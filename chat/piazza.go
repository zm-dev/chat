package chat

import (
	"encoding/json"
	"time"
)

type IMsg interface {
	GetUser() User
	GetData() []byte
	GetSendAt() time.Time
	GetMeta() map[string]string
	SetMeta(map[string]string)
	MarshalJSON() ([]byte, error)
}

type Msg struct {
	User   User
	Data   []byte
	Meta   map[string]string
	SendAt time.Time
}

func (m *Msg) GetUser() User {
	if m == nil {
		return nil
	}
	return m.User
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

func (m *Msg) GetMeta() map[string]string {
	return m.Meta
}

func (m *Msg) SetMeta(meta map[string]string) {
	m.Meta = meta
}

func (m *Msg) MarshalJSON() ([]byte, error) {
	m.GetUser().GetId()
	return json.Marshal(struct {
		UserId int64             `json:"user_id"`
		Data   string            `json:"data"`
		Meta   map[string]string `json:"meta"`
		SendAt time.Time         `json:"send_at"`
	}{
		UserId: m.User.GetId(),
		Data:   string(m.Data),
		Meta:   m.Meta,
		SendAt: m.SendAt,
	})
}

type SysUser struct {
}

func (SysUser) GetId() int64 {
	return 0
}

func (SysUser) SendMsg(msg IMsg) error {
	return nil
}

type Piazza struct {
	chatRooms map[string]*ChatRoom
}

func (p *Piazza) CreateChatRoom(Id string) *ChatRoom {
	cr := NewChatRoom(Id)
	p.chatRooms[cr.Id] = cr
	go cr.Run()
	return cr
}

func (p *Piazza) Receive(chatRoomId string, m IMsg) {
	if cr, ok := p.chatRooms[chatRoomId]; ok {
		cr.Broadcast(m)
	}
}

func (p *Piazza) Join(chatRoomId string, user User) {
	if cr, ok := p.chatRooms[chatRoomId]; ok {
		cr.Join(user)
	}
}

func (p *Piazza) Leave(chatRoomId string, user User) {
	if cr, ok := p.chatRooms[chatRoomId]; ok {
		cr.Leave(user)
	}
}

func NewPiazza() *Piazza {
	return &Piazza{chatRooms: make(map[string]*ChatRoom)}
}
