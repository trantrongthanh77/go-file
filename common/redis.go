package common

import (
	"context"
	"go-file/common/config"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var RedisEnabled = true

// InitRedisClient This function is called after init()
func InitRedisClient(conf *config.Config) (err error) {
	if os.Getenv("REDIS_CONN_STRING") == "" {
		RedisEnabled = false
		// The cache depends on Redis
		ExplorerCacheEnabled = false
		// This stat feature also depends on Redis
		StatEnabled = false
		return nil
	}
	opt, err := redis.ParseURL(conf.RedisConnectionString)
	if err != nil {
		panic(err)
	}
	RDB = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RDB.Ping(ctx).Result()
	return err
}

func ParseRedisOption(conf *config.Config) *redis.Options {
	opt, err := redis.ParseURL(conf.RedisConnectionString)
	if err != nil {
		panic(err)
	}
	return opt
}
