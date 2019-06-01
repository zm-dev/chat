package store

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/zm-dev/chat_v2/model"
	"github.com/zm-dev/chat_v2/store/db_store"
	"github.com/zm-dev/chat_v2/store/redis_store"
)

type Store interface {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.RecordStore
}

type store struct {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.RecordStore
}

func NewStore(db *gorm.DB, redisClient *redis.Client) Store {
	return &store{redis_store.NewRedisTicket(redisClient),
		db_store.NewDBUser(db),
		db_store.NewDBCertificate(db),
		db_store.NewDBRecord(db),
	}
}
