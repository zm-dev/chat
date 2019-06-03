package service

import (
	"context"
	"github.com/zm-dev/chat/model"
)

type ChatRoomService struct {
	model.ChatRoomStore
}

func (cSvc *ChatRoomService) ChatRoomLoad(id string) (*model.ChatRoom, error) {
	return cSvc.ChatRoomStore.ChatRoomLoad(id)
}

func (cSvc *ChatRoomService) ChatRoomCreate(cr *model.ChatRoom) error {
	return cSvc.ChatRoomStore.ChatRoomCreate(cr)
}

func (cSvc *ChatRoomService) ChatRoomDelete(id string) error {
	return cSvc.ChatRoomStore.ChatRoomDelete(id)
}

func (cSvc *ChatRoomService) ChatRoomList() []*model.ChatRoom {
	return cSvc.ChatRoomStore.ChatRoomList()
}

func ChatRoomLoad(c context.Context, id string) (*model.ChatRoom, error) {
	return FromContext(c).ChatRoomLoad(id)
}

func ChatRoomCreate(c context.Context, cr *model.ChatRoom) error {
	return FromContext(c).ChatRoomCreate(cr)
}

func ChatRoomDelete(c context.Context, id string) error {
	return FromContext(c).ChatRoomDelete(id)
}

func ChatRoomList(c context.Context) []*model.ChatRoom {
	return FromContext(c).ChatRoomList()
}

func NewChatRoom(cs model.ChatRoomStore) model.ChatRoomService {
	return &ChatRoomService{cs}
}
