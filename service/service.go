package service

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/zm-dev/chat/config"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/pkg/hasher"
	"github.com/zm-dev/chat/store"
	"time"
)

type Service interface {
	model.TicketService
	model.UserService
	model.CertificateService
	model.RecordService
	model.ChatService
}

type service struct {
	model.TicketService
	model.UserService
	model.CertificateService
	model.RecordService
	model.ChatService
}

func NewService(db *gorm.DB, redisClient *redis.Client, baseFs afero.Fs, conf *config.Config) Service {
	s := store.NewStore(db, redisClient)
	tSvc := NewTicketService(s, time.Duration(conf.Ticket.TTL)*time.Second)
	//h := hasher.NewBcyptHasher(
	//	[]byte(conf.AppSalt),
	//	3,
	//	32<<10,
	//	uint8(runtime.NumCPU()),
	//	32,
	//)
	h := hasher.NewBcyptHasher()
	return &service{
		tSvc,
		NewUserService(s, s, tSvc, h),
		NewCertificateService(s),
		NewRecordService(s),
		NewChatService(),
	}
}
