package server

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/zm-dev/chat_v2/config"
	"github.com/zm-dev/chat_v2/pkg/pubsub"
	"github.com/zm-dev/chat_v2/service"
	"go.uber.org/zap"
)

type Server struct {
	AppEnv      string
	Debug       bool
	BaseFs      afero.Fs
	RedisClient *redis.Client
	DB          *gorm.DB
	Conf        *config.Config
	Logger      *zap.Logger
	Service     service.Service
	Pub         pubsub.PubQueue
}
