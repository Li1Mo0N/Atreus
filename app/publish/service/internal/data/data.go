package data

import (
	"Atreus/app/publish/service/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewPublishRepo, NewMysqlConn, NewRedisConn)

// Data .
type Data struct {
	db    *gorm.DB
	cache *redis.Client
	log   *log.Helper
}

// NewData .
func NewData(db *gorm.DB, cacheClient *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	data := &Data{
		db:    db,
		cache: cacheClient,
		log:   log.NewHelper(logger),
	}
	return data, cleanup, nil
}

// NewMysqlConn mysql数据库连接
func NewMysqlConn(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Mysql.Dsn))
	if err != nil {
		log.Fatalf("Database connection failure, err : %v", err)
	}
	InitDB(db)
	log.Info("Database enabled successfully!")
	return db
}

// NewRedisConn Redis数据库连接
func NewRedisConn(c *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:           int(c.Redis.Db),
		Addr:         c.Redis.Addr,
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		Password:     c.Redis.Password,
	})

	// ping Redis客户端，判断连接是否存在
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis database connection failure, err : %v", err)
	}
	log.Info("Cache enabled successfully!")
	return client
}

func InitDB(db *gorm.DB) {
	if err := db.AutoMigrate(&Video{}); err != nil {
		log.Fatalf("Database initialization error, err : %v", err)
	}
}
