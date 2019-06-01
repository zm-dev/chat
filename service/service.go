package service

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/zm-dev/chat_v2/config"
	"github.com/zm-dev/chat_v2/model"
	"github.com/zm-dev/chat_v2/pkg/hasher"
	"github.com/zm-dev/chat_v2/pkg/pubsub"
	"github.com/zm-dev/chat_v2/store"
	"runtime"
	"time"
)

type Service interface {
	model.TicketService
	model.UserService
	model.CertificateService
}

type service struct {
	model.TicketService
	model.UserService
	model.CertificateService
}

func NewService(db *gorm.DB, redisClient *redis.Client, baseFs afero.Fs, conf *config.Config, pub pubsub.PubQueue) Service {
	s := store.NewStore(db, redisClient)
	tSvc := NewTicketService(s, time.Duration(conf.Ticket.TTL)*time.Second)
	h := hasher.NewArgon2Hasher(
		[]byte(conf.AppSalt),
		3,
		32<<10,
		uint8(runtime.NumCPU()),
		32,
	)
	return &service{
		tSvc,
		NewUserService(s, s, tSvc, h),
		NewCertificateService(s),
	}
}