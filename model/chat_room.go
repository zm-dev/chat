package model

import "errors"

type ChatRoom struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CoverHash string `json:"cover_hash"`
}

type ChatRoomStore interface {
	ChatRoomLoad(string) (*ChatRoom, error)
	ChatRoomCreate(*ChatRoom) error
	ChatRoomDelete(string) error
	ChatRoomList() []*ChatRoom
}

type ChatRoomService interface {
	ChatRoomStore
}

var (
	ErrChatRoomNotExist = errors.New("chat_room 不存在")
)
