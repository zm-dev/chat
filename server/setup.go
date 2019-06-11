package server

import (
	"fmt"
	nosConfig "github.com/NetEase-Object-Storage/nos-golang-sdk/config"
	"github.com/NetEase-Object-Storage/nos-golang-sdk/nosclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/afero"
	"github.com/wq1019/go-image_uploader"
	"github.com/wq1019/go-image_uploader/image_url"
	imageUploaderNos "github.com/wq1019/go-image_uploader/nos"
	"github.com/zm-dev/chat/config"
	"github.com/zm-dev/chat/model"
	"github.com/zm-dev/chat/service"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func setupGorm(debug bool, databaseConfig *config.DatabaseConfig) *gorm.DB {
	var dataSourceName string
	switch databaseConfig.Driver {
	case "sqlite3":
		dataSourceName = databaseConfig.DBName
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			databaseConfig.User,
			databaseConfig.Password,
			databaseConfig.Host+":"+databaseConfig.Port,
			databaseConfig.DBName,
		)
	}

	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(databaseConfig.Driver, dataSourceName)
		if err == nil {
			db.LogMode(debug)
			if debug {
				// ONLY_FULL_GROUP_BY,
				db.Exec("set session sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'")
				autoMigrate(db)
			}
			return db
		}
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！ error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Certificate{},
		&model.Record{},
		&image_uploader.Image{},
		&model.LastRecord{},
	).Error
	if err != nil {
		log.Fatalf("AutoMigrate 失败！ error: %+v", err)
	}
}

func setupRedis(redisConfig *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address + ":" + redisConfig.Port,
		Password: redisConfig.Password,
	})
}

func setupNos(s *Server) *nosclient.NosClient {
	nosClient, err := nosclient.New(&nosConfig.Config{
		Endpoint:  s.Conf.Nos.Endpoint,
		AccessKey: s.Conf.Nos.AccessKey,
		SecretKey: s.Conf.Nos.SecretKey,
	})
	if err != nil {
		log.Fatalf("nos client 创建失败! error: %+v", err)
	}
	return nosClient
}

func setupImageUploader(s *Server) image_uploader.Uploader {
	nosClient := setupNos(s)
	return imageUploaderNos.NewNosUploader(
		image_uploader.HashFunc(image_uploader.MD5HashFunc),
		image_uploader.NewDBStore(s.DB),
		nosClient,
		s.Conf.Nos.BucketName,
		image_uploader.Hash2StorageNameFunc(image_uploader.TwoCharsPrefixHash2StorageNameFunc),
	)
}

func setupImageURL(s *Server) image_url.URL {
	return image_url.NewNosImageProxyURL(
		s.Conf.ImageProxy.Host,
		s.Conf.Nos.ExternalEndpoint,
		s.Conf.Nos.BucketName,
		s.Conf.ImageProxy.OmitBaseUrl == "true",
		image_uploader.Hash2StorageNameFunc(image_uploader.TwoCharsPrefixHash2StorageNameFunc),
	)
}

func setupFilesystem(fsConfig *config.FilesystemConfig) afero.Fs {
	switch fsConfig.Driver {
	case "os":
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	case "memory":
		return afero.NewBasePathFs(afero.NewMemMapFs(), fsConfig.Root)
	default:
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	}
}

func loadEnv(appEnv string) string {
	if appEnv == "" {
		appEnv = "production"
	}
	return appEnv
}

func setupLogger(srv *Server) *zap.Logger {
	var err error
	var logger *zap.Logger
	if srv.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func SetupServer(configPath string) *Server {
	s := &Server{}

	s.AppEnv = loadEnv(os.Getenv("APP_ENV"))

	s.Debug = os.Getenv("DEBUG") == "true"

	s.Logger = setupLogger(s)

	s.Logger.Debug("load config...")
	s.Conf = config.LoadConfig(configPath)

	s.Logger.Debug("load filesystem...")
	s.BaseFs = setupFilesystem(&s.Conf.Fs)

	s.Logger.Debug("load redis...")
	s.RedisClient = setupRedis(&s.Conf.Redis)

	s.Logger.Debug("load database...")
	s.DB = setupGorm(s.Debug, &s.Conf.DB)

	s.Logger.Debug("load service...")
	s.Service = service.NewService(s.DB, s.RedisClient, s.BaseFs, s.Conf)

	s.Logger.Debug("load uploader service...")
	s.ImageUploader = setupImageUploader(s)

	s.ImageUrl = setupImageURL(s)

	s.NosClient = setupNos(s)
	return s
}
