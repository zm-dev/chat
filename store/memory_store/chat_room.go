package memory_store

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
	"github.com/zm-dev/chat/model"
	"sync"
)

type chatRoom struct {
	chatRooms sync.Map
}

func (c *chatRoom) ChatRoomLoad(id string) (*model.ChatRoom, error) {
	chatRoom, ok := c.chatRooms.Load(id)
	if !ok {
		return nil, model.ErrChatRoomNotExist
	}
	return chatRoom.(*model.ChatRoom), nil
}

func (c *chatRoom) ChatRoomList() []*model.ChatRoom {
	chatRooms := make([]*model.ChatRoom, 0, 3)
	c.chatRooms.Range(func(key, v interface{}) bool {
		chatRooms = append(chatRooms, v.(*model.ChatRoom))
		return true
	})
	return chatRooms
}

func (c *chatRoom) ChatRoomCreate(chatRoom *model.ChatRoom) error {
	u4 := uuid.NewV4()

	cid := hex.EncodeToString(u4.Bytes())
	chatRoom.Id = cid
	c.chatRooms.Store(cid, chatRoom)
	return nil
}

func (c *chatRoom) ChatRoomDelete(id string) error {
	if _, ok := c.chatRooms.Load(id); !ok {
		return model.ErrChatRoomNotExist
	}
	c.chatRooms.Delete(id)
	return nil
}
func NewMemoryChatRoom() model.ChatRoomStore{
	return &chatRoom{
		sync.Map{},
	}
}
