package store

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/store/db_store"
	"github.com/zm-dev/chat/store/memory_store"
)

type Store interface {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.RecordStore
	model.ChatRoomStore
}

type store struct {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.RecordStore
	model.ChatRoomStore
}

func NewStore(db *gorm.DB, redisClient *redis.Client) Store {
	return &store{memory_store.NewMemoryTicket(),
		db_store.NewDBUser(db),
		db_store.NewDBCertificate(db),
		db_store.NewDBRecord(db),
		memory_store.NewMemoryChatRoom(),
	}
}
