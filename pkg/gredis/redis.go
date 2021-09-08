package gredis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/lin07ux/go-gin-example/pkg/setting"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func Setup() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         setting.RedisSetting.Host,
		Password:     setting.RedisSetting.Password,
		DB:           setting.RedisSetting.Index,
		PoolSize:     setting.RedisSetting.PoolSize,
		IdleTimeout:  setting.RedisSetting.IdleTimeout,
		MinIdleConns: 1,
	})
}

func Set(key string, data interface{}, ttl time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, ttl).Err()
}

func Get(key string, v interface{}) error {
	reply, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(reply, v)
}

func Exists(key string) (bool, error) {
	val, err := rdb.Exists(ctx, key).Result()

	return val > 0, err
}

func Delete(key string) (bool, error) {
	val, err := rdb.Del(ctx, key).Result()

	return val > 0, err
}

func LikeDeletes(key string) error {
	iter := rdb.Scan(ctx, 0, "*" + key + "*", 0).Iterator()

	for iter.Next(ctx) {
		if err := iter.Err(); err != nil {
			return err
		}

		if _, err := Delete(iter.Val()); err != nil {
			return err
		}
	}

	return nil
}
