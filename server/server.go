package server

import (
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/wq1019/go-image_uploader"
	"github.com/wq1019/go-image_uploader/image_url"
	"github.com/zm-dev/chat/config"
	"github.com/zm-dev/chat/pkg/pubsub"
	"github.com/zm-dev/chat/service"
	"go.uber.org/zap"
)

type Server struct {
	AppEnv        string
	Debug         bool
	BaseFs        afero.Fs
	RedisClient   *redis.Client
	DB            *gorm.DB
	Conf          *config.Config
	Logger        *zap.Logger
	ImageUploader image_uploader.Uploader
	ImageUrl      image_url.URL
	NosClient     *nosclient.NosClient
	Service       service.Service
	Pub           pubsub.PubQueue
}
